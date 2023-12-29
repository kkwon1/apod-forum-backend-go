package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/db/dao"
	"github.com/kkwon1/apod-forum-backend/internal/models"
	"github.com/kkwon1/apod-forum-backend/internal/models/requests"
	"github.com/kkwon1/apod-forum-backend/internal/repositories"
)

var dbClient *db.MongoDBClient
var apodRepository *repositories.ApodRepository
var userRepository *repositories.UserRepository

var apodDao *dao.ApodDao
var postUpvoteDao *dao.PostUpvoteDao

var astroTerms map[string]struct{}

func main() {
	initialize()
	startService()
}

func initialize() {
	loadEnvFile()
	mongoConnectionString := os.Getenv("MONGO_ENDPOINT")

	dbClient, err := db.NewMongoDBClient(mongoConnectionString)
	if err != nil {
		log.Fatal("Error connecting to Mongo DB")
	}

	apodDao, _ = dao.NewApodDao(dbClient)
	postUpvoteDao, _ = dao.NewPostUpvoteDao(dbClient)
	apodRepository, _ = repositories.NewApodRepository(apodDao)
	userRepository, _ = repositories.NewUserRepository(postUpvoteDao)
}

func loadEnvFile() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	loadErr := godotenv.Load(curDir + "/.env")
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir)
	}
}

func startService() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowOrigins = []string{os.Getenv("ALLOWED_ORIGINS")}
	r.Use(cors.New(config))

	verifyJwt := getJwtVerifierMiddleware()

	// APOD
	r.GET("/apods/random", getRandomApod)
	r.GET("/apods/:date", getApod)
	r.GET("/apods", getApodPage)
	r.GET("/apods/random/:count", getRandomApods)
	r.GET("/apods/search", searchApod)

	// Posts
	r.GET("/posts/:id", getPost)
	r.POST("/posts/upvote", upvote)

	// Commments
	r.GET("/comments/:id", getComment)

	// Users
	r.GET("/users/:userSub", verifyJwt, getUser)

	// Load the Astronomy terms
	astroTerms = loadAstroTerms()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// TODO: Verify claims and make sure you only allow the correct user
func getJwtVerifierMiddleware() gin.HandlerFunc {
	issuerURL, _ := url.Parse(os.Getenv("JWT_ISSUER"))
	audience := os.Getenv("AUTH0_AUDIENCE")

	provider := jwks.NewCachingProvider(issuerURL, time.Duration(5*time.Minute))

	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
	)

	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	return adapter.Wrap(jwtMiddleware.CheckJWT)
}

// ========== APOD ==========

func getRandomApod(c *gin.Context) {
	resultChan := make(chan models.Apod)

	go func() {
		apod := apodRepository.GetRandomApod()
		apod.Tags = extractTags(apod, astroTerms)
		resultChan <- apod
	}()

	apod := <-resultChan
	c.JSON(http.StatusOK, apod)
}

func getRandomApods(c *gin.Context) {
	// TODO implement
}

func getApod(c *gin.Context) {
	date := c.Param("date")
	apod := apodRepository.GetApod(date)
	apod.Tags = extractTags(apod, astroTerms)
	c.JSON(http.StatusOK, apod)
}

func getApodPage(c *gin.Context) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	today := time.Now()
	endDate := today.AddDate(0, 0, (-1 * offset))
	startDate := endDate.AddDate(0, 0, (-1 * (limit - 1)))

	apods := apodRepository.GetApodsBetweenDates(startDate, endDate)
	c.JSON(http.StatusOK, apods)
}

func searchApod(c *gin.Context) {
	searchString := c.Query("searchString")
	apods := apodRepository.SearchApods(searchString)

	c.JSON(http.StatusOK, apods)
}

// ==================== POSTS ==================

func getPost(c *gin.Context) {
	date := c.Param("id")

	// Create a channel to receive the result
	resultChan := make(chan models.ApodPost)
	go func() {
		post := apodRepository.GetApodPost(date)
		post.NasaApod.Tags = extractTags(post.NasaApod, astroTerms)
		resultChan <- post
	}()
	post := <-resultChan
	c.JSON(http.StatusOK, post)
}

func upvote(c *gin.Context) {
	var upvoteRequest requests.UpvotePostRequest

	if err := c.BindJSON(&upvoteRequest); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, "")
	}

	userRepository.UpvotePost(upvoteRequest)
	apodRepository.IncrementUpvoteCount(upvoteRequest.PostId)
}

// ========================= Comments ===========================

func getComment(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world")
}

func getUser(c *gin.Context) {
	userSub := c.Param("userSub")
	postIds := userRepository.GetUpvotedPostIds(userSub)

	var user models.User
	user = models.User{
		UserSub:           userSub,
		UserName:          "testUsername",
		Email:             "testEmail",
		EmailVerified:     true,
		ProfilePictureUrl: "testProfileUrl",
		UpvotedPostIds:    postIds,
	}
	c.JSON(http.StatusOK, user)
}


func loadAstroTerms() (map[string]struct{}) {
	file, err := os.Open("internal/const/astro_terms.txt")
    if err != nil {
        log.Fatalf("failed to open file: %s", err)
    }
    defer file.Close()

    set := make(map[string]struct{})
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        set[strings.ToLower(scanner.Text())] = struct{}{}
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("failed to scan file: %s", err)
    }

    // Print the set
    return set
}

func extractTags(apod models.Apod, astro_terms map[string]struct{}) []string {
	words := strings.Fields(strings.ToLower(apod.Explanation))
	matches := make(map[string]struct{})
	for _, word := range words {
			if _, ok := astro_terms[word]; ok {
					matches[word] = struct{}{}
			}
	}
	return setToList(matches)
}

func setToList(set map[string]struct{}) []string {
	list := make([]string, 0, len(set))
	for key := range set {
			list = append(list, key)
	}
	return list
}