package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kkwon1/apod-forum-backend/cmd/models"
	"github.com/kkwon1/apod-forum-backend/cmd/repositories"
	"github.com/kkwon1/apod-forum-backend/cmd/utils"
)

type ApodController struct {
	router *gin.Engine
	apodRepository *repositories.ApodRepository
}

func NewApodController(router *gin.Engine, apodRepository *repositories.ApodRepository) (*ApodController, error) {
	return &ApodController{router: router, apodRepository: apodRepository}, nil
}

func (ac *ApodController) RegisterRoutes() {
	apodRoute := ac.router.Group("/apods")
	apodRoute.GET("/random", ac.getRandomApod)
	apodRoute.GET("/:date", ac.getApod)
	apodRoute.GET("", ac.getApodPage)
	apodRoute.GET("/random/:count", ac.getRandomApods)
	apodRoute.GET("/search", ac.searchApod)
}

func (ac *ApodController) getRandomApod(c *gin.Context) {
	resultChan := make(chan models.Apod)

	go func() {
		apod := ac.apodRepository.GetRandomApod()
		apod.Tags = utils.ExtractTags(apod)
		resultChan <- apod
	}()

	apod := <-resultChan
	c.JSON(http.StatusOK, apod)
}

func (ac *ApodController) getRandomApods(c *gin.Context) {
	// TODO implement
}

func (ac *ApodController) getApod(c *gin.Context) {
	date := c.Param("date")
	apod := ac.apodRepository.GetApod(date)
	apod.Tags = utils.ExtractTags(apod)
	c.JSON(http.StatusOK, apod)
}

func (ac *ApodController) getApodPage(c *gin.Context) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	today := time.Now()
	endDate := today.AddDate(0, 0, (-1 * offset))
	startDate := endDate.AddDate(0, 0, (-1 * (limit - 1)))

	apods := ac.apodRepository.GetApodsBetweenDates(startDate, endDate)
	c.JSON(http.StatusOK, apods)
}

func (ac *ApodController) searchApod(c *gin.Context) {
	searchString := c.Query("searchString")
	apods := ac.apodRepository.SearchApods(searchString)

	c.JSON(http.StatusOK, apods)
}