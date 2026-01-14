package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var urlsMap map[string]string = make(map[string]string)

func main() {
	router := gin.Default()
	router.POST("/create", createShortLink)
	router.GET("/get", expendShortLink)
	router.Run("localhost:8080")
}

func createShortLink(ginContext *gin.Context) {
	urlToShorten := ginContext.GetString("url")
	isUrl := validateUrl(urlToShorten)
	if isUrl {
		shortUrl := generateShortUrl(urlToShorten)
	}
	ginContext.String(http.StatusOK)
}

func generateShortUrl(urlToShorten string) string {
	"localhost:8080/" + ""
}

func validateUrl(urlToShorten string) bool {
	_, err := url.Parse(urlToShorten)
	if err != nil {
		return false
	}

	return true
}

func expendShortLink(ginContext *gin.Context) {

}
