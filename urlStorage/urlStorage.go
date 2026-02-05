package urlStorage

// UrlStorage is an interface for URL storage implementations
type UrlStorage interface {
	// Store saves the mapping between the original URL and the shortened URL
	Store(string string)
	Get(string) string
}
