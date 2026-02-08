package urlStorage

type InMemoryUrlStorage struct {
	realToShortMap map[string]string
	shortToRealMap map[string]string
}

func (inMemoryStorage *InMemoryUrlStorage) StoreStore(urlToShorten string, shortUrl string) {
	inMemoryStorage.realToShortMap[urlToShorten] = shortUrl
}

func (inMemoryStorage *InMemoryUrlStorage) Get(shortUrl string) string {
	return inMemoryStorage.shortToRealMap[shortUrl]
}
