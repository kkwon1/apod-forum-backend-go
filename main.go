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
	"github.com/kkwon1/apod-forum-backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *db.MongoDBClient
var apodCache *lru.Cache[string, Apod]

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
	if apodCache.Contains(date) {
		var apod, _ = apodCache.Get(date)
		c.JSON(http.StatusOK, apod)
	} else {
		apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")
		var apod Apod
		filter := bson.M{"date": date}
		apodCollection.FindOne(context.Background(), filter).Decode(&apod)

		apodCache.Add(date, apod)

		c.JSON(http.StatusOK, apod)
	}
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
	NasaApod Apod `json:"nasaApod"`
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
		}

		c.JSON(http.StatusOK, post)
	}
}
