package main

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortUrl(serverUrl string, urlSize int, customAlias string) (string, error) {

	path, error := GetPath(customAlias, urlSize)
	if error != nil {
		return "", error
	}

	hostUrl, _ := url.Parse(serverUrl)
	shortUrl := hostUrl.JoinPath(path)
	return shortUrl.String(), nil
}

func GenerateRandomPath(urlSize int) string {
	var randomPart string = ""
	for range urlSize {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		randBase62Char := base62[randIndex.Int64()]
		randomPart += string(randBase62Char)
	}
	return randomPart
}

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
