package urlStorage

// UrlPair represents a pair of short and original URLs
type UrlPair struct {
	// ShortUrl is the shortened URL
	ShortUrl string `json:"short_url"`
	// OriginalUrl is the original URL that was shortened
	OriginalUrl string `json:"original_url"`
}
