package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kkwon1/apod-forum-backend/cmd/controllers"
	"github.com/kkwon1/apod-forum-backend/cmd/db"
	"github.com/kkwon1/apod-forum-backend/cmd/db/dao"
	"github.com/kkwon1/apod-forum-backend/cmd/repositories"
)

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

func startService() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowOrigins = []string{os.Getenv("ALLOWED_ORIGINS")}
	r.Use(cors.New(config))

	r.GET("/", getComment)

	// APOD
	apodController, _ := controllers.NewApodController(r, apodRepository)
	userController, _ := controllers.NewUserController(r, userRepository)
	postController, _ := controllers.NewPostController(r, apodRepository, userRepository)

	apodController.RegisterRoutes()
	userController.RegisterRoutes()
	postController.RegisterRoutes();

	r.GET("/comments/:id", getComment)

	r.Run()
}

func loadEnvFile() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	envFilePath := curDir + "/.env"
	_, err = os.Stat(envFilePath)
	if os.IsNotExist(err) {
		log.Println(".env file does not exist. Using environment variables on host")
		return
	}

	loadErr := godotenv.Load(envFilePath)
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir)
	}
}


func getComment(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world")
}