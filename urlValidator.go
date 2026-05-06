package main

import (
	"errors"
	"log/slog"
	"net/url"
	"unicode"
	"unicode/utf8"
)

const (
	minAliasLenght = 2
	maxAliasLenght = 20
	testhostName   = "http://localhost:8080"
)

func validateUrl(urlToShorten string) bool {
	_, err := url.ParseRequestURI(urlToShorten)
	if err != nil {
		slog.Warn("Invalid URL", "url", urlToShorten)
		return false
	}

	return true
}

// Validate if custom
func ValidateAlias(aliasToValidate string) error {
	if aliasToValidate == "" {
		return errors.New("Alias is empty")
	}

	runesCount := utf8.RuneCountInString(aliasToValidate)
	if (runesCount < minAliasLenght) || (runesCount > maxAliasLenght) {
		return errors.New("Invalid alias size")
	}

	for _, aliasRune := range aliasToValidate {
		runeValidationError := validateAliasRune(aliasRune)
		if runeValidationError != nil {
			return runeValidationError
		}
	}

	_, urlParseError := url.Parse(testhostName)
	if urlParseError != nil {
		slog.Error("Failed to custom alias test url", "error", urlParseError)
		return errors.New("Unsupportead alias " + aliasToValidate)
	}
	return nil
}

func validateAliasRune(aliasRune rune) error {
	if unicode.IsControl(aliasRune) {
		return errors.New("Alias contains control characters")
	}

	if unicode.IsSpace(aliasRune) {
		return errors.New("Alias contains whitespace characters")
	}

	if unicode.IsPunct(aliasRune) {
		return errors.New("Alias contains punctuation characters")
	}

	if aliasRune == '/' || aliasRune == '?' || aliasRune == '#' || aliasRune == '&' || aliasRune == '=' {
		errorText := "Alias contains reserved URL characters: " + string(aliasRune)
		return errors.New(errorText)
	}

	return nil
}
