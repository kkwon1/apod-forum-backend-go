package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kkwon1/apod-forum-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

var dbClient *db.MongoDBClient

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

	apodCollection := dbClient.GetDatabase("apodDB").Collection("apod")
	var result Apod
	filter := bson.M{"date": "2023-01-22"}
	apodCollection.FindOne(context.Background(), filter).Decode(&result)

	log.Println(result)
	// get_apod_from_nasa()
	start_service()
}

func start_service() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

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

type NasaApodDAO struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// We don't actually ever need this functionality, because we have a lambda that is running once a day
// to pull the APOD data and store in mongodb
func get_apod_from_nasa() {
	api_key := os.Getenv("NASA_API_KEY")
	resp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", api_key))
	if err != nil {
		log.Fatal("Failed to call NASA API")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var apod NasaApodDAO
	if err := json.Unmarshal(body, &apod); err != nil {
		log.Fatal("Cannot unmarshal JSON")
	}

	fmt.Println(PrettyPrint(apod))
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
