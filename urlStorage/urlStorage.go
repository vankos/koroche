package urlStorage

// UrlStorage is an interface for URL storage implementations
type UrlStorage interface {
	// Saves the mapping between the original URL and the shortened URL
	Store(string string)
	// Retrieves the original URL based on the shortened URL
	Get(string) string
}
