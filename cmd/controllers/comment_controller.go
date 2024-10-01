package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
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