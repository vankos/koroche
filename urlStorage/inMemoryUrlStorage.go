package urlStorage

// InMemoryUrlStorage is an in-memory implementation of the UrlStorage interface
type InMemoryUrlStorage struct {
	realToShortMap map[string]string
	shortToRealMap map[string]string
	clickCountMap  map[string]int
}

// Saves the mapping between the original URL and the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) Store(urlToShorten string, shortUrl string) error {
	inMemoryStorage.shortToRealMap[shortUrl] = urlToShorten
	inMemoryStorage.realToShortMap[urlToShorten] = shortUrl
	return nil
}

// Retrieves the original URL based on the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetOriginalUrl(shortUrl string) (string, error) {
	return inMemoryStorage.shortToRealMap[shortUrl], nil
}

// Increments the click count for a given shortened URL
func (inMemoryStorage *InMemoryUrlStorage) IncrementClick(shortUrl string) error {
	inMemoryStorage.clickCountMap[shortUrl]++
	return nil
}

// Retrieves the click count for a given shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetClickCount(shortUrl string) (int, error) {
	return inMemoryStorage.clickCountMap[shortUrl], nil
}

// // Gets saved short URL for a given original URL, "" if not found
func (inMemoryStorage *InMemoryUrlStorage) GetShortUrl(url string) (string, error) {
	return inMemoryStorage.realToShortMap[url], nil
}
