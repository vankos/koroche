package main

import (
	"context"
	"io"
	"koroche/urlStorage"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller is responsible for handling incoming HTTP requests and interacting with the URL storage
type Controller struct {
	statsChannel  chan string
	urlStorage    urlStorage.UrlStorage
	shortUrlsHost string
}

func NewController(urlStorage urlStorage.UrlStorage, statsChannel chan string, shortUrlsHost string) *Controller {
	controller := &Controller{
		urlStorage:    urlStorage,
		statsChannel:  statsChannel,
		shortUrlsHost: shortUrlsHost,
	}

	go controller.CollectStats()
	return controller
}

func (controller *Controller) ProcessCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := validateUrl(urlToShorten)
	if !isUrl {
		slog.Info("Invalid URL provided for shortening", "url", urlToShorten)
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx := ginContext.Request.Context()
	existingShortUrl, _ := controller.urlStorage.GetShortUrl(ctx, urlToShorten)
	if existingShortUrl != "" {
		ginContext.String(http.StatusOK, existingShortUrl)
		return
	}

	customAlias := ginContext.Query("custom_alias")
	shortUrl, aliasError := GenerateShortUrl(controller.shortUrlsHost, shortUrlSize, customAlias)
	if aliasError != nil {
		slog.Info("Invalid custom alias provided for shortening", "alias", customAlias, "error", aliasError)
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	err := controller.urlStorage.Store(ctx, urlToShorten, shortUrl)
	if err != nil {
		slog.Info("	Failed to store URL mapping", "url", urlToShorten, "error", err)
		ginContext.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ginContext.String(http.StatusOK, shortUrl)
}

func (controller *Controller) ProcessGet(ginContext *gin.Context) {
	urlToExpand := ginContext.Query("url")
	isUrl := validateUrl(urlToExpand)
	if !isUrl {
		slog.Info("	Invalid URL provided for expansion", "url", urlToExpand)
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx := ginContext.Request.Context()
	realUrl, ok := controller.urlStorage.GetOriginalUrl(ctx, urlToExpand)
	if ok != nil {
		slog.Info("	No original URL found for the provided short URL", "url", urlToExpand)
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
		slog.Info("	Invalid URL provided for stats retrieval", "url", url)
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx := ginContext.Request.Context()
	stats, ok := controller.urlStorage.GetStats(ctx, url)
	if ok != nil {
		slog.Info("	No stats found for the provided URL", "url", url)
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.JSON(http.StatusOK, stats)
}

func (controller *Controller) CollectStats() {
	slog.Info("Starting stats collector")
	backgrpundContext := context.Background()
	for shortUrl := range controller.statsChannel {
		func() {
			ctx, cancel := context.WithTimeout(backgrpundContext, time.Second*5)
			defer cancel()
			controller.urlStorage.UpdateStats(ctx, shortUrl)
		}()
	}
}

func (controller *Controller) Close() error {
	slog.Info("Closing controller")
	closable, ok := controller.urlStorage.(io.Closer)
	if ok {
		closable.Close()
	}

	return nil
}
