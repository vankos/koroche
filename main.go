package main

import (
	"context"
	"io"
	"koroche/urlStorage"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const serverUrl = "http://localhost:8080"
const shortUrlSize = 12

func main() {
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage()
	defer CloseIfNeeded(&urlStorage)
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan string)
	controller := NewController(&urlStorage, statsChannel)
	go CollectStats(statsChannel, &urlStorage)
	router := gin.Default()
	router.POST("/create", controller.ProcessCreate)
	router.GET("/get", controller.ProcessGet)
	router.GET("/stats", controller.ProcessStats)
	parsedServerUrl, _ := url.ParseRequestURI(serverUrl)
	router.Run(parsedServerUrl.Host)
	ctx := context.Background()
	go MonitorOldLinks(ctx, &urlStorage)
}

func CollectStats(statsChannel <-chan string, urlStorage urlStorage.UrlStorage) {
	for shortUrl := range statsChannel {
		urlStorage.UpdateStats(shortUrl)
	}
}

func CloseIfNeeded(urlStorage urlStorage.UrlStorage) {
	closable, ok := urlStorage.(io.Closer)
	if !ok {
		return
	}

	closable.Close()
}

func validateUrl(urlToShorten string) bool {
	_, err := url.ParseRequestURI(urlToShorten)
	if err != nil {
		return false
	}

	return true
}

func MonitorOldLinks(ctx context.Context, urlStorage urlStorage.UrlStorage) {
	DeleteOldLinks(urlStorage)
	ticker := time.NewTicker(time.Hour * 24)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			DeleteOldLinks(urlStorage)
		}
	}
}

func DeleteOldLinks(urlStorage urlStorage.UrlStorage) {
	urlStorage.DeleteLinksOlderThan(time.Hour * 24)
}
