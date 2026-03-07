package main

import (
	"context"
	"io"
	"koroche/urlStorage"
	"net/http"
	"time"

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

	go controller.CollectStats()
	return controller
}

func (controller *Controller) ProcessCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := validateUrl(urlToShorten)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	ctx := ginContext.Request.Context()
	existingShortUrl, _ := controller.urlStorage.GetShortUrl(ctx, urlToShorten)
	if existingShortUrl != "" {
		ginContext.String(http.StatusOK, existingShortUrl)
		return
	}

	shortUrl := GenerateShortUrl(serverUrl, shortUrlSize)
	controller.urlStorage.Store(ctx, urlToShorten, shortUrl)
	ginContext.String(http.StatusOK, shortUrl)
}

func (controller *Controller) ProcessGet(ginContext *gin.Context) {
	urlToExpand := ginContext.Query("url")
	isUrl := validateUrl(urlToExpand)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx := ginContext.Request.Context()
	realUrl, ok := controller.urlStorage.GetOriginalUrl(ctx, urlToExpand)
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

	ctx := ginContext.Request.Context()
	stats, ok := controller.urlStorage.GetStats(ctx, url)
	if ok != nil {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.JSON(http.StatusOK, stats)
}

func (controller *Controller) CollectStats() {
	backgrpundContext := context.Background()
	for shortUrl := range controller.statsChannel {
		ctx, cancel := context.WithTimeout(backgrpundContext, time.Second*5)
		defer cancel()
		controller.urlStorage.UpdateStats(ctx, shortUrl)
	}
}

func (controller *Controller) Close() error {
	closable, ok := controller.urlStorage.(io.Closer)
	if ok {
		closable.Close()
	}

	close(controller.statsChannel)
	return nil
}
