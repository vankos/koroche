package controller

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/vankos/koroche/urlStorage"
	"github.com/vankos/koroche/urlTools"

	"github.com/gin-gonic/gin"
)

// Controller is responsible for handling incoming HTTP requests and interacting with the URL storage
type Controller struct {
	statsChannel  chan string
	urlStorage    urlStorage.UrlStorage
	shortUrlsHost string
	shortUrlSize  int
}

// NewController creates a new instance of Controller,
// initializes the stats collector goroutine, and returns the controller instance
func NewController(
	urlStorage urlStorage.UrlStorage,
	statsChannel chan string,
	shortUrlsHost string,
	shortUrlSize int) *Controller {
	controller := &Controller{
		urlStorage:    urlStorage,
		statsChannel:  statsChannel,
		shortUrlsHost: shortUrlsHost,
		shortUrlSize:  shortUrlSize,
	}

	go controller.CollectStats()
	return controller
}

// ProcessCreate handles the creation of a shortened URL.
// It validates the input URL, checks for existing mappings,
// generates a new short URL if necessary, and stores the mapping in the URL storage.
func (controller *Controller) ProcessCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := urlTools.ValidateUrl(urlToShorten)
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
	shortUrl, aliasError := urlTools.GenerateShortUrl(controller.shortUrlsHost, controller.shortUrlSize, customAlias)
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

// ProcessGet handles the retrieval of the original URL based on the provided shortened URL.
func (controller *Controller) ProcessGet(ginContext *gin.Context) {
	urlToExpand := ginContext.Query("url")
	isUrl := urlTools.ValidateUrl(urlToExpand)
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

// ProcessStats handles the retrieval of statistics for a given shortened URL, including click count and original URL.
func (controller *Controller) ProcessStats(ginContext *gin.Context) {
	url := ginContext.Query("url")
	isUrl := urlTools.ValidateUrl(url)
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

// CollectStats listens for short URLs on the statsChannel and updates their statistics in the URL storage.
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

// SetupRouter initializes the Gin router and registers the controller's endpoints for creating,
// retrieving, and getting statistics of shortened URLs.
func SetupRouter(controller *Controller) *gin.Engine {
	router := gin.Default()
	router.POST("/create", controller.ProcessCreate)
	router.GET("/get", controller.ProcessGet)
	router.GET("/stats", controller.ProcessStats)
	return router
}

// Close gracefully shuts down the controller, ensuring that any resources used by the URL storage are properly released.
func (controller *Controller) Close() error {
	slog.Info("Closing controller")
	closable, ok := controller.urlStorage.(io.Closer)
	if ok {
		closable.Close()
	}

	return nil
}
