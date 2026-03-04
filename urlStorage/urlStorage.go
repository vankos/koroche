package urlStorage

import (
	"context"
	"time"
)

// UrlStorage is an interface for URL storage implementations
type UrlStorage interface {
	// Saves the mapping between the original URL and the shortened URL
	Store(context.Context, string, string) error
	// Retrieves the original URL based on the shortened URL
	GetOriginalUrl(context.Context, string) (string, error)
	// Retrieves the link stats for a given shortened URL
	GetStats(context.Context, string) (LinkStats, error)
	// Gets saved short URL for a given original URL, "" if not found
	GetShortUrl(context.Context, string) (string, error)
	// Deletes links that are older than the specified duration
	DeleteLinksOlderThan(context.Context, time.Duration)
	// Update stats for short URL
	UpdateStats(context.Context, string) error
}
