package shortUrlGenerator

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortUrl(serverUrl string, urlSize int) string {

	var randomPart string = ""
	for range urlSize {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		randBase62Char := base62[randIndex.Int64()]
		randomPart += string(randBase62Char)
	}

	hostUrl, _ := url.Parse(serverUrl)
	shortUrl := hostUrl.JoinPath(randomPart)
	return shortUrl.String()
}
