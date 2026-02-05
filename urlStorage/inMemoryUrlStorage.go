package urlStorage

type InMemeoryUrlStorage struct {
	realToShortMap map[string]string
	shortToRealMap map[string]string
}

func (inMemoryStorage *InMemeoryUrlStorage) StoreStore(urlToShorten string, shortUrl string) {
	inMemoryStorage.realToShortMap[urlToShorten] = shortUrl
}
