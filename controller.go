package main

import (
	"koroche/urlStorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller is responsible for handling incoming HTTP requests and interacting with the URL storage
type Controller struct {
	statsChannel chan string
	urlStorage   urlStorage.UrlStorage
}

func NewController(urlStorage urlStorage.UrlStorage, statsChannel chan string) *Controller {
	controller := &Controller{
		urlStorage:   urlStorage,
		statsChannel: statsChannel,
	}

	return controller
}

func (controller *Controller) ProcessCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := validateUrl(urlToShorten)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	existingShortUrl, _ := controller.urlStorage.GetShortUrl(urlToShorten)
	if existingShortUrl != "" {
		ginContext.String(http.StatusOK, existingShortUrl)
		return
	}

	shortUrl := GenerateShortUrl(serverUrl, shortUrlSize)
	controller.urlStorage.Store(urlToShorten, shortUrl)
	ginContext.String(http.StatusOK, shortUrl)
}

func (controller *Controller) ProcessGet(ginContext *gin.Context) {
	urlToExpand := ginContext.Query("url")
	isUrl := validateUrl(urlToExpand)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	realUrl, ok := controller.urlStorage.GetOriginalUrl(urlToExpand)
	if ok != nil {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.String(http.StatusOK, realUrl)
	controller.statsChannel <- urlToExpand
}

func (controller *Controller) ProcessStats(ginContext *gin.Context) {
	url := ginContext.Query("url")
	isUrl := validateUrl(url)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	stats, ok := controller.urlStorage.GetStats(url)
	if ok != nil {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.JSON(http.StatusOK, stats)
}
