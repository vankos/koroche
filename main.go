package main

import (
	"koroche/shortUrlGenerator"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const host = "localhost:8080"
const shortUrlSize = 12

var urlsMap map[string]string = make(map[string]string)

func main() {
	router := gin.Default()
	router.POST("/create", processCreate)
	router.GET("/get", processGet)
	router.Run(host)
}

func processCreate(ginContext *gin.Context) {
	urlToShorten := ginContext.GetString("url")
	isUrl := validateUrl(urlToShorten)
	if isUrl {
		shortUrl := generateShortUrl(urlToShorten)
		ginContext.String(http.StatusOK, shortUrl)
	}

	ginContext.AbortWithStatus(http.StatusUnprocessableEntity)
}

func generateShortUrl(urlToShorten string) string {
	shortUrl := shortUrlGenerator.GenerateShortUrl(host, shortUrlSize)
	urlsMap[urlToShorten] = shortUrl
	return shortUrl
}

func validateUrl(urlToShorten string) bool {
	_, err := url.Parse(urlToShorten)
	if err != nil {
		return false
	}

	return true
}

func processGet(ginContext *gin.Context) {

}
