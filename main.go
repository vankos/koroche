package main

import (
	"context"
	"koroche/urlStorage"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const shortUrlSize = 12

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))
	slog.SetDefault(logger)
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage()
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan string, 100)
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
	slog.Info("Closing..")
	controller.Close()
}

func validateUrl(urlToShorten string) bool {
	_, err := url.ParseRequestURI(urlToShorten)
	if err != nil {
		slog.Warn("Invalid URL", "url", urlToShorten)
		return false
	}

	return true
}

func MonitorOldLinks(ctx context.Context, urlStorage urlStorage.UrlStorage) {
	slog.Info("Starting old links monitor")
	DeleteOldLinks(ctx, urlStorage)
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping old links monitor")
			return
		case <-ticker.C:
			DeleteOldLinks(ctx, urlStorage)
		}
	}
}

func DeleteOldLinks(ctx context.Context, urlStorage urlStorage.UrlStorage) {
	urlStorage.DeleteLinksOlderThan(ctx, time.Hour*24)
}
