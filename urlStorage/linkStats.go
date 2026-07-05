package urlStorage

import "time"

// LinkStats represents the statistics of a shortened URL
type LinkStats struct {
	// ShortUrl is the shortened URL
	ShortUrl string
	// OriginalUrl is the original URL that was shortened
	OriginalUrl string
	// CreatedAt is the time when the short URL was created
	CreatedAt time.Time
	// LastAccessedAt is the time when the short URL was last accessed
	LastAccessedAt *time.Time
	// ClickCount is the number of times the short URL has been accessed
	ClickCount int
}
