package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/joho/godotenv"
	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *db.MongoDBClient
var apodCache *lru.Cache[string, Apod]
var apodRepository *repositories.ApodRepository

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoConnectionString := os.Getenv("MONGO_ENDPOINT")

	dbClient, err = db.NewMongoDBClient(mongoConnectionString)
	if err != nil {
		log.Fatal("Error connecting to Mongo DB")
	}

	apodRepository, _ = repositories.NewApodRepository(dbClient)

	// initialize LRU Cache with 3000 items
	apodCache, _ = lru.New[string, Apod](3000)

	start_service()
}

func start_service() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	// APOD
	r.GET("/apod/random", getRandomApod)
	r.GET("/apod/:date", getApod)
	r.GET("/apods", getApodPage)
	r.GET("/apods/:count", getRandomApods)
	r.GET("/apods/search", searchApod)

	// Posts
	r.GET("/posts/:id", getPost)

	// Commments
	r.GET("/comments/:id", getComment)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// ========== APOD ==========
type Apod struct {
	Copyright      string `json:"copyright"`
	PostID         string `json:"postId"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"mediaType"`
	ServiceVersion string `json:"serviceVersion"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Hdurl          string `json:"hdurl"`
	UpvoteCount    int    `json:"upvoteCount"`
	SaveCount      int    `json:"saveCount"`
	CommentCount   int    `json:"commentCount"`
}

func getRandomApod(c *gin.Context) {
	apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")

	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.M{"size": 1}}},
	}
	cursor, _ := apodCollection.Aggregate(context.Background(), pipeline)

	defer cursor.Close(context.Background())

	// Check if there are results
	if cursor.Next(context.Background()) {
		var apod Apod
		err := cursor.Decode(&apod)
		if err != nil {
			log.Fatal(err)
		}

		apodCache.Add(apod.Date, apod)
		c.JSON(http.StatusOK, apod)
	} else {
		fmt.Println("No random document found.")
	}
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

	apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")
	filter := bson.M{
		"date": bson.M{
			"$gte": startDate.Format("2006-01-02"),
			"$lte": endDate.Format("2006-01-02"),
		},
	}
	cursor, err := apodCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Failed to read apod page")
	}
	defer cursor.Close(context.Background())

	// Process the results
	var results []Apod
	for cursor.Next(context.Background()) {
		var apod Apod
		if err := cursor.Decode(&apod); err != nil {
			log.Fatal("Failed to decode APOD")
		}
		results = append(results, apod)
	}

	c.JSON(http.StatusOK, results)
}

func searchApod(c *gin.Context) {
	searchString := c.Query("searchString")
	apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")
	pipeline := mongo.Pipeline{
		{{
			Key: "$search", Value: bson.M{
				"index": "textSearch",
				"text": bson.M{
					"query": searchString,
					"path": bson.M{
						"wildcard": "*",
					},
				},
			},
		}},
	}
	cursor, _ := apodCollection.Aggregate(context.Background(), pipeline)

	defer cursor.Close(context.Background())

	// Process the results
	var results []Apod
	for cursor.Next(context.Background()) {
		var apod Apod
		if err := cursor.Decode(&apod); err != nil {
			log.Fatal("Failed to decode APOD")
		}
		results = append(results, apod)
	}

	c.JSON(http.StatusOK, results)
}

// ==================== POSTS ==================
type ApodPost struct {
	NasaApod Apod        `json:"nasaApod"`
	Comments CommentTree `json:"comments"`
}

func getPost(c *gin.Context) {
	date := c.Param("id")
	if apodCache.Contains(date) {
		var apod, _ = apodCache.Get(date)

		var post ApodPost = ApodPost{
			NasaApod: apod,
		}

		c.JSON(http.StatusOK, post)
	} else {
		apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")
		var apod Apod
		filter := bson.M{"date": date}
		apodCollection.FindOne(context.Background(), filter).Decode(&apod)
		apodCache.Add(date, apod)

		var post ApodPost = ApodPost{
			NasaApod: apod,
			Comments: CommentTree{
				CommentID:    date,
				Children:     []CommentTree{},
				CreateDate:   "2023-09-30",
				ModifiedDate: "2023-09-30",
				Comment:      "Test Comment",
				Author:       "Test Author",
				IsDeleted:    false,
				IsLeaf:       false,
			},
		}

		c.JSON(http.StatusOK, post)
	}
}

// ========================= Comments ===========================

/*
   private static final String COMMENT_PATH = "/comment";

   @Autowired
   private CommentsClient commentsClient;

   @Autowired
   private ApodClient apodClient;

   @PostMapping(path = COMMENT_PATH + "/add", consumes = MediaType.APPLICATION_JSON_VALUE)
   public Comment addComment(@RequestBody AddCommentRequest addCommentRequest) {
       Comment comment = commentsClient.addComment(addCommentRequest);
       apodClient.addCommentToPost(addCommentRequest.getPostId());
       return comment;
   }

   @PostMapping(path = COMMENT_PATH + "/delete", consumes = MediaType.APPLICATION_JSON_VALUE, produces = MediaType.APPLICATION_JSON_VALUE)
   public String deleteComment(@RequestBody DeleteCommentRequest deleteCommentRequest) {
       return commentsClient.softDeleteComment(deleteCommentRequest);
   }

   @GetMapping(path = COMMENT_PATH, params = {"post_id"})
   public CommentTree getPostComments(@RequestParam String post_id) {
       return commentsClient.getPostComments(post_id);
   }

	 public class CommentTree {
    String commentId;
    List<CommentTree> children;
    LocalDateTime createDate;
    LocalDateTime modifiedDate;
    String comment;
    String author;
    Boolean isDeleted;
    Boolean isLeaf;
}
*/

type CommentTree struct {
	CommentID    string        `json:"commentId"`
	Children     []CommentTree `json:"children"`
	CreateDate   string        `json:"createDate"`
	ModifiedDate string        `json:"modifiedDate"`
	Comment      string        `json:"comment"`
	Author       string        `json:"author"`
	IsDeleted    bool          `json:"isDeleted"`
	IsLeaf       bool          `json:"isLeaf"`
}

func getComment(c *gin.Context) {

	c.JSON(http.StatusOK, "hello world")
}
