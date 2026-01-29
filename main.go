package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const serverUrl = "http://localhost:8080"
const shortUrlSize = 12

var realToShortMap map[string]string = make(map[string]string)
var shortToRealMap map[string]string = make(map[string]string)

func main() {
	router := gin.Default()
	router.POST("/create", processCreate)
	router.GET("/get", processGet)
	parsedServerUrl, _ := url.ParseRequestURI(serverUrl)
	router.Run(parsedServerUrl.Host)
}

func processCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.Query("url")
	isUrl := validateUrl(urlToShorten)
	if isUrl {
		shortUrl := generateShortUrl(urlToShorten)
		ginContext.String(http.StatusOK, shortUrl)
	}

	ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
}

func generateShortUrl(urlToShorten string) string {
	shortUrl := GenerateShortUrl(serverUrl, shortUrlSize)
	realToShortMap[urlToShorten] = shortUrl
	shortToRealMap[shortUrl] = urlToShorten
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

	realUrl, ok := shortToRealMap[urlToExpand]
	if !ok {
		ginContext.AbortWithStatus(http.StatusNotFound)
		return
	}

	ginContext.String(http.StatusOK, realUrl)
}
