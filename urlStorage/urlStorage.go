package urlStorage

// UrlStorage is an interface for URL storage implementations
type UrlStorage interface {
	// Saves the mapping between the original URL and the shortened URL
	Store(string, string) error
	// Retrieves the original URL based on the shortened URL
	GetOriginalUrl(string) (string, error)
	// Increments the click count for a given shortened URL
	IncrementClick(string) error
	// Retrieves the link stats for a given shortened URL
	GetStats(string) (LinkStats, error)
	// Gets saved short URL for a given original URL, "" if not found
	GetShortUrl(string) (string, error)
}
