package main

import (
	"context"
	"github.com/vankos/koroche/urlStorage"
	"log/slog"
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
	dsn := os.Getenv("POSTGRES_DATA_SOURCE_NAME")
	urlStorage, err := urlStorage.NewPostgreSqlUrlStorage(dsn)
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan string, 100)
	shortUrlHostName := os.Getenv("HOST_NAME")
	controller := NewController(&urlStorage, statsChannel, shortUrlHostName)
	defer Close(controller)
	router := SetupRouter(controller)
	port := ":" + os.Getenv("PORT")
	bgCtx := context.Background()
	ctx, cancel := context.WithCancel(bgCtx)
	defer cancel()
	go MonitorOldLinks(ctx, &urlStorage)
	router.Run(port)

}

func SetupRouter(controller *Controller) *gin.Engine {
	router := gin.Default()
	router.POST("/create", controller.ProcessCreate)
	router.GET("/get", controller.ProcessGet)
	router.GET("/stats", controller.ProcessStats)
	return router
}

func Close(controller *Controller) {
	slog.Info("Closing..")
	controller.Close()
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
