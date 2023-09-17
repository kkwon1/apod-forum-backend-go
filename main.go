package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	get_apod_from_nasa()
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

	fmt.Println(string(body))
}