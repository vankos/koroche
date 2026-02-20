package urlStorage

import "time"

// InMemoryUrlStorage is an in-memory implementation of the UrlStorage interface
type InMemoryUrlStorage struct {
	realToShortMap map[string]string
	shortToRealMap map[string]*LinkStats
}

// Saves the mapping between the original URL and the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) Store(urlToShorten string, shortUrl string) error {
	inMemoryStorage.shortToRealMap[shortUrl] = &LinkStats{
		ShortUrl:     shortUrl,
		OriginalUrl:  urlToShorten,
		CreatedAt:    time.Now(),
		LasAccesedAt: time.Now(),
		ClickCount:   0}

	inMemoryStorage.realToShortMap[urlToShorten] = shortUrl
	return nil
}

// Retrieves the original URL based on the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetOriginalUrl(shortUrl string) (string, error) {
	inMemoryStorage.shortToRealMap[shortUrl].ClickCount++
	inMemoryStorage.shortToRealMap[shortUrl].LasAccesedAt = time.Now()
	return inMemoryStorage.shortToRealMap[shortUrl].OriginalUrl, nil
}

// Retrieves the click count for a given shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetStats(shortUrl string) (LinkStats, error) {
	return *inMemoryStorage.shortToRealMap[shortUrl], nil
}

// // Gets saved short URL for a given original URL, "" if not found
func (inMemoryStorage *InMemoryUrlStorage) GetShortUrl(url string) (string, error) {
	return inMemoryStorage.realToShortMap[url], nil
}
