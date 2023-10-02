package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kkwon1/apod-forum-backend/internal/db"
	"github.com/kkwon1/apod-forum-backend/internal/db/dao"
	"github.com/kkwon1/apod-forum-backend/internal/repositories"
)

var dbClient *db.MongoDBClient
var apodRepository *repositories.ApodRepository

var apodDao *dao.ApodDao

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

	apodDao, _ = dao.NewApodDao(dbClient)
	apodRepository, _ = repositories.NewApodRepository(apodDao)

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
