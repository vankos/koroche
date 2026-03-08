package urlStorage

import "time"

type LinkStats struct {
	ShortUrl     string
	OriginalUrl  string
	CreatedAt    time.Time
	LasAccesedAt *time.Time
	ClickCount   int
}
