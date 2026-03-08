package main

import (
	"context"
	"koroche/urlStorage"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const serverUrl = "http://localhost:8080"
const shortUrlSize = 12

func main() {
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage()
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan string)
	controller := NewController(&urlStorage, statsChannel)
	defer Close(*controller)
	router := gin.Default()
	router.POST("/create", controller.ProcessCreate)
	router.GET("/get", controller.ProcessGet)
	router.GET("/stats", controller.ProcessStats)
	parsedServerUrl, _ := url.ParseRequestURI(serverUrl)
	bgCtx := context.Background()
	ctx, cancel := context.WithCancel(bgCtx)
	defer cancel()
	go MonitorOldLinks(ctx, &urlStorage)
	router.Run(parsedServerUrl.Host)

}

func Close(controller Controller) {
	controller.Close()
}

func validateUrl(urlToShorten string) bool {
	_, err := url.ParseRequestURI(urlToShorten)
	if err != nil {
		return false
	}

	return true
}

func MonitorOldLinks(ctx context.Context, urlStorage urlStorage.UrlStorage) {
	DeleteOldLinks(ctx, urlStorage)
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			DeleteOldLinks(ctx, urlStorage)
		}
	}
}

func DeleteOldLinks(ctx context.Context, urlStorage urlStorage.UrlStorage) {
	urlStorage.DeleteLinksOlderThan(ctx, time.Hour*24)
}
