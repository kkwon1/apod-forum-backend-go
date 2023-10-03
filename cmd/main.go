package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	"github.com/kkwon1/apod-forum-backend/internal/repositories"
)

var dbClient *db.MongoDBClient
var apodRepository *repositories.ApodRepository
var userRepository *repositories.UserRepository

var apodDao *dao.ApodDao
var postUpvoteDao *dao.PostUpvoteDao

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
	config.AllowOrigins = []string{os.Getenv("ALLOWED_ORIGINS")}
	r.Use(cors.New(config))

	verifyJwt := getJwtVerifierMiddleware()

	// APOD
	r.GET("/apods/random", verifyJwt, getRandomApod)
	r.GET("/apods/:date", getApod)
	r.GET("/apods", getApodPage)
	r.GET("/apods/random/:count", getRandomApods)
	r.GET("/apods/search", searchApod)

	// Posts
	r.GET("/posts/:id", getPost)

	// Commments
	r.GET("/comments/:id", getComment)

	// Users
	r.GET("/users/:userSub", getUser)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

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
	apod := apodRepository.GetRandomApod()
	c.JSON(http.StatusOK, apod)
}

func getRandomApods(c *gin.Context) {
	// TODO implement
}

func getApod(c *gin.Context) {
	date := c.Param("date")
	apod := apodRepository.GetApod(date)
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
	post := apodRepository.GetApodPost(date)
	c.JSON(http.StatusOK, post)
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
