package urlStorage

import "time"

// LinkStats represents the statistics of a shortened URL
type LinkStats struct {
	UrlPair
	// CreatedAt is the time when the short URL was created
	CreatedAt time.Time `json:"created_at"`
	// LastAccessedAt is the time when the short URL was last accessed
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
	// ClickCount is the number of times the short URL has been accessed
	ClickCount int `json:"click_count"`
}
