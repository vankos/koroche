package urlStorage

import "time"

type LinkStats struct {
	ShortUrl      string
	OriginalUrl   string
	CreatedAt     time.Time
	LastAccesedAt *time.Time
	ClickCount    int
}
