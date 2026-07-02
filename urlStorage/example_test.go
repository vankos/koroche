package urlStorage_test

import (
	"context"
	"fmt"
	"time"

	"github.com/vankos/koroche/urlStorage"
)

func ExampleUrlStorage_Store() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	err := urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	fmt.Println(err == nil)
	// Output: true
}

func ExampleUrlStorage_GetOriginalUrl() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	origianUrl, _ := urlsStorage.GetOriginalUrl(ctx, "http://short.url/abc123")
	fmt.Println(origianUrl)
	// Output: https://www.google.com
}

func ExampleUrlStorage_UpdateStats() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	// Click  x2
	urlsStorage.UpdateStats(ctx, "http://short.url/abc123")
	urlsStorage.UpdateStats(ctx, "http://short.url/abc123")
	stats, _ := urlsStorage.GetStats(ctx, "http://short.url/abc123")
	fmt.Println(stats.ClickCount)
	// Output: 2
}

func ExampleUrlStorage_GetStats() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	urlsStorage.UpdateStats(ctx, "http://short.url/abc123")
	stats, _ := urlsStorage.GetStats(ctx, "http://short.url/abc123")
	fmt.Println(stats.ClickCount)
	fmt.Println(stats.OriginalUrl)
	fmt.Println(stats.ShortUrl)
	// Output: 1
	// https://www.google.com
	// http://short.url/abc123
}

func ExampleUrlStorage_GetShortUrl() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	shortUrl, _ := urlsStorage.GetShortUrl(ctx, "https://www.google.com")
	fmt.Println(shortUrl)
	// Output: http://short.url/abc123
}

func ExampleUrlStorage_DeleteLinksOlderThan() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	urlsStorage.Store(ctx, "https://www.google.com", "http://short.url/abc123")
	time.Sleep(10 * time.Millisecond)
	urlsStorage.DeleteLinksOlderThan(ctx, 5*time.Millisecond)
	_, err := urlsStorage.GetShortUrl(ctx, "https://www.google.com")
	fmt.Println(err != nil)
	// Output: true
}

func ExampleUrlStorage_GetTopLinks() {
	urlsStorage := urlStorage.NewInMemoryUrlStorage()
	ctx := context.Background()
	originalUrl1 := "https://www.example.com"
	shortUrl1 := "abc"
	originalUrl2 := "https://www.example1.com"
	shortUrl2 := "abc111"
	originalUrl3 := "https://www.example2.com"
	shortUrl3 := "abc222"
	urlsStorage.Store(ctx, originalUrl1, shortUrl1)
	urlsStorage.Store(ctx, originalUrl2, shortUrl2)
	urlsStorage.Store(ctx, originalUrl3, shortUrl3)
	urlsStorage.UpdateStats(ctx, shortUrl2)
	urlsStorage.UpdateStats(ctx, shortUrl2)
	urlsStorage.UpdateStats(ctx, shortUrl2)
	urlsStorage.UpdateStats(ctx, shortUrl1)
	urlsStorage.UpdateStats(ctx, shortUrl1)
	urlsStorage.UpdateStats(ctx, shortUrl3)
	toplinls, _ := urlsStorage.GetTopLinks(ctx, 2, 1)
	fmt.Println(len(toplinls) == 2)
	fmt.Println(toplinls[0].ShortUrl)
	fmt.Println(toplinls[1].ShortUrl)
	// Output: true
	// abc
	// abc222
}
