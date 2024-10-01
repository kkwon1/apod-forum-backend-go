package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"github.com/kkwon1/apod-forum-backend/cmd/models/requests"
	"github.com/kkwon1/apod-forum-backend/cmd/repositories"
)

type CommentController struct {
	router *gin.Engine
	apodRepository *repositories.ApodRepository
}

func NewCommentController(router *gin.Engine, apodRepository *repositories.ApodRepository) (*CommentController, error) {
	return &CommentController{router: router, apodRepository: apodRepository}, nil
}

func (cc *CommentController) RegisterRoutes() {
	commentRoute := cc.router.Group("/comments")
	commentRoute.GET("/:postId", cc.getCommentsForPost)
	commentRoute.POST("/:postId", cc.addComment)
}

func (cc *CommentController) getCommentsForPost(c *gin.Context) {
	resultChan := make(chan []*models.CommentNode)
	postId := c.Param("postId")

	go func() {
		comments := cc.apodRepository.GetCommentsForPost(postId)
		resultChan <- comments
	}()

	commentsResult := <-resultChan
	c.JSON(http.StatusOK, commentsResult)
}

func (cc *CommentController) addComment(c *gin.Context) {
	var addCommentRequest requests.AddCommentRequest

	if err := c.BindJSON(&addCommentRequest); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, "")
	}

	comment := &models.Comment{
		PostID:    addCommentRequest.PostID,
		ParentID:  addCommentRequest.ParentID,
		CommentID:   uuid.New().String(),
		Comment:  addCommentRequest.Comment,
		Author:  addCommentRequest.Author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cc.apodRepository.AddCommentForPost(*comment)
}

