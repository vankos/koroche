package urlStorage

import (
	"cmp"
	"context"
	"slices"
	"time"
)

// InMemoryUrlStorage is an in-memory implementation of the UrlStorage interface
type InMemoryUrlStorage struct {
	realToShortMap map[string]string
	shortToRealMap map[string]*LinkStats
}

// NewInMemoryUrlStorage ceates a new instance of InMemoryUrlStorage
func NewInMemoryUrlStorage() InMemoryUrlStorage {
	return InMemoryUrlStorage{
		realToShortMap: make(map[string]string),
		shortToRealMap: make(map[string]*LinkStats),
	}
}

// Saves the mapping between the original URL and the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) Store(_ context.Context, urlToShorten string, shortUrl string) error {
	now := time.Now()
	inMemoryStorage.shortToRealMap[shortUrl] = &LinkStats{
		ShortUrl:      shortUrl,
		OriginalUrl:   urlToShorten,
		CreatedAt:     now,
		LastAccesedAt: &now,
		ClickCount:    0}

	inMemoryStorage.realToShortMap[urlToShorten] = shortUrl
	return nil
}

// Retrieves the original URL based on the shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetOriginalUrl(_ context.Context, shortUrl string) (string, error) {
	stats, ok := inMemoryStorage.shortToRealMap[shortUrl]
	if !ok {
		return "", &UrlStorageError{Code: NotFound}
	}

	return stats.OriginalUrl, nil
}

// Retrieves the click count for a given shortened URL
func (inMemoryStorage *InMemoryUrlStorage) GetStats(_ context.Context, shortUrl string) (LinkStats, error) {
	stats, ok := inMemoryStorage.shortToRealMap[shortUrl]
	if !ok {
		return LinkStats{}, &UrlStorageError{Code: NotFound}
	}

	return *stats, nil
}

// Gets saved short URL for a given original URL, "" if not found
func (inMemoryStorage *InMemoryUrlStorage) GetShortUrl(_ context.Context, url string) (string, error) {
	result, ok := inMemoryStorage.realToShortMap[url]
	if !ok {
		return "", &UrlStorageError{Code: NotFound}
	}

	return result, nil
}

// Deletes links that are older than the specified duration
func (inMemoryStorage *InMemoryUrlStorage) DeleteLinksOlderThan(_ context.Context, olderThan time.Duration) {
	for key := range inMemoryStorage.shortToRealMap {
		stats := inMemoryStorage.shortToRealMap[key]
		oldTime := time.Now().Add(-olderThan)
		if stats.CreatedAt.Before(oldTime) {
			delete(inMemoryStorage.shortToRealMap, key)
			delete(inMemoryStorage.realToShortMap, stats.OriginalUrl)
		}
	}
}

// Update stats for short URL
func (inMemoryStorage *InMemoryUrlStorage) UpdateStats(_ context.Context, shortUrl string) error {
	inMemoryStorage.shortToRealMap[shortUrl].ClickCount++
	now := time.Now()
	inMemoryStorage.shortToRealMap[shortUrl].LastAccesedAt = &now
	return nil
}

func (inMemoryStorage *InMemoryUrlStorage) GetTopLinks(ctx context.Context, urlsToReturn int, offset int) ([]string, error) {
	mapValues := make([]LinkStats, 0, len(inMemoryStorage.shortToRealMap))
	for _, value := range inMemoryStorage.shortToRealMap {
		mapValues = append(mapValues, *value)
	}
	slices.SortFunc(mapValues, func(a, b LinkStats) int {
		return cmp.Compare(a.ClickCount, b.ClickCount)
	})

	topDoamins := make([]string, 0, urlsToReturn)
	for i := offset; i < offset+urlsToReturn; i++ {
		topDoamins = append(topDoamins, mapValues[i].OriginalUrl)
	}

	return topDoamins, nil

}
