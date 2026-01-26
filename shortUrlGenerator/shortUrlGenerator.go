package shortUrlGenerator

import (
	"crypto/rand"
	"math/big"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortUrl(host string, urlSize int) string {

	var randomPart string = ""
	for range urlSize {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		randBase62Char := base62[randIndex.Int64()]
		randomPart += string(randBase62Char)
	}

	shortUrl := host
	return shortUrl
}
