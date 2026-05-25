package main

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateShortUrl generates a short URL based on the provided server URL, desired URL size, and an optional custom alias.
func GenerateShortUrl(serverUrl string, urlSize int, customAlias string) (string, error) {

	path, error := GetPath(customAlias, urlSize)
	if error != nil {
		return "", error
	}

	hostUrl, _ := url.Parse(serverUrl)
	shortUrl := hostUrl.JoinPath(path)
	return shortUrl.String(), nil
}

// GenerateRandomPath generates a random string of the specified length using characters from the base62 set.
func GenerateRandomPath(urlSize int) string {
	var randomPart string = ""
	for range urlSize {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		randBase62Char := base62[randIndex.Int64()]
		randomPart += string(randBase62Char)
	}
	return randomPart
}

// GetPath determines the path to be used for the short URL based on the provided custom alias and desired URL size.
func GetPath(customAlias string, shortUrlSize int) (string, error) {
	if customAlias == "" {
		return GenerateRandomPath(shortUrlSize), nil
	}

	validationError := ValidateAlias(customAlias)
	if validationError != nil {
		return "", validationError
	}

	return customAlias, nil
}
