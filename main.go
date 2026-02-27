package main

import (
	"context"
	"io"
	"koroche/urlStorage"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const serverUrl = "http://localhost:8080"
const shortUrlSize = 12

var urlStorageObj urlStorage.UrlStorage

func main() {
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage()
	urlStorageObj = &urlStorage
	defer CloseIfNeeded(urlStorageObj)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.POST("/create", processCreate)
	router.GET("/get", processGet)
	router.GET("/stats", processStats)
	parsedServerUrl, _ := url.ParseRequestURI(serverUrl)
	router.Run(parsedServerUrl.Host)
	ctx := context.Background()
	go MonitorOldLinks(ctx)
}

func CloseIfNeeded(urlStorage urlStorage.UrlStorage) {
	closable, ok := urlStorage.(io.Closer)
	if !ok {
		return
	}

	closable.Close()
}

func processCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := validateUrl(urlToShorten)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
	}

	existingShortUrl, _ := urlStorageObj.GetShortUrl(urlToShorten)
	if existingShortUrl != "" {
		ginContext.String(http.StatusOK, existingShortUrl)
		return
	}

	shortUrl := generateShortUrl(urlToShorten)
	urlStorageObj.Store(urlToShorten, shortUrl)
	ginContext.String(http.StatusOK, shortUrl)
}

func generateShortUrl(urlToShorten string) string {
	shortUrl := GenerateShortUrl(serverUrl, shortUrlSize)
	urlStorageObj.Store(urlToShorten, shortUrl)
	return shortUrl
}

func validateUrl(urlToShorten string) bool {
	_, err := url.ParseRequestURI(urlToShorten)
	if err != nil {
		return false
	}

	return true
}

func processGet(ginContext *gin.Context) {
	urlToExpand := ginContext.Query("url")
	isUrl := validateUrl(urlToExpand)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	realUrl, ok := urlStorageObj.GetOriginalUrl(urlToExpand)
	if ok != nil {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.String(http.StatusOK, realUrl)
}

func processStats(ginContext *gin.Context) {
	url := ginContext.Query("url")
	isUrl := validateUrl(url)
	if !isUrl {
		ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	stats, ok := urlStorageObj.GetStats(url)
	if ok != nil {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.JSON(http.StatusOK, stats)
}

func MonitorOldLinks(ctx context.Context) {
	DeleteOldLinks()
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			DeleteOldLinks()
		}
	}
}

func DeleteOldLinks() {
	urlStorageObj.DeleteLinksOlderThan(time.Hour * 24)
}
