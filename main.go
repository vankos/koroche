package main

import (
	"context"
	"koroche/urlStorage"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const shortUrlSize = 12

func main() {
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage()
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan string)
	shortUrlHostName := os.Getenv("HOST_NAME")
	controller := NewController(&urlStorage, statsChannel, shortUrlHostName)
	defer Close(controller)
	router := gin.Default()
	router.POST("/create", controller.ProcessCreate)
	router.GET("/get", controller.ProcessGet)
	router.GET("/stats", controller.ProcessStats)
	port := ":" + os.Getenv("PORT")
	bgCtx := context.Background()
	ctx, cancel := context.WithCancel(bgCtx)
	defer cancel()
	go MonitorOldLinks(ctx, &urlStorage)
	router.Run(port)

}

func Close(controller *Controller) {
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
