package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"github.com/kkwon1/apod-forum-backend/cmd/models/requests"
	"github.com/kkwon1/apod-forum-backend/cmd/repositories"
	"github.com/kkwon1/apod-forum-backend/cmd/utils"
)

type PostController struct {
	router *gin.Engine
	apodRepository *repositories.ApodRepository
	userRepository *repositories.UserRepository
}

func NewPostController(router *gin.Engine, apodRepository *repositories.ApodRepository, userRepository *repositories.UserRepository) (*PostController, error) {
	return &PostController{router: router, apodRepository: apodRepository, userRepository: userRepository}, nil
}

func (pc *PostController) RegisterRoutes() {
	postGroup := pc.router.Group("/posts")
	postGroup.GET("/:id", pc.getPost)
	postGroup.POST("/upvote", pc.upvote)
}

func (pc *PostController) getPost(c *gin.Context) {
	date := c.Param("id")

	// Create a channel to receive the result
	resultChan := make(chan models.ApodPost)
	go func() {
		post := pc.apodRepository.GetApodPost(date)
		post.NasaApod.Tags = utils.ExtractTags(post.NasaApod)
		resultChan <- post
	}()
	post := <-resultChan
	c.JSON(http.StatusOK, post)
}

func (pc *PostController) upvote(c *gin.Context) {
	var upvoteRequest requests.UpvotePostRequest

	if err := c.BindJSON(&upvoteRequest); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, "")
	}

	pc.userRepository.UpvotePost(upvoteRequest)
	pc.apodRepository.IncrementUpvoteCount(upvoteRequest.PostId)
}


