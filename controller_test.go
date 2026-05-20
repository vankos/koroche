package main

import (
	"encoding/json"
	"fmt"
	"github.com/vankos/koroche/urlStorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	urlStorage := urlStorage.NewInMemoryUrlStorage()
	statsChannel := make(chan string, 100)
	shortUrlHostName := "http://localhost:8080"
	controller := NewController(&urlStorage, statsChannel, shortUrlHostName)
	router := SetupRouter(controller)
	return router
}

func TestСreate(t *testing.T) {
	router := SetUpRouter()
	testData := []struct {
		requestUrl   string
		expectedCode int
		expectString bool
	}{
		{"https://www.google.com", http.StatusOK, true},
		{"invalid-url", http.StatusUnprocessableEntity, false},
	}

	for _, data := range testData {
		t.Run(data.requestUrl, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/create?url="+data.requestUrl, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			assert.Equal(t, data.expectedCode, recorder.Code)
			assert.Equal(t, recorder.Body.String() != "", data.expectString)
		})

	}
}

func TestGet(t *testing.T) {
	router := SetUpRouter()
	originalUrl := "https://www.google.com"
	request, _ := http.NewRequest("POST", "/create?url="+originalUrl, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	shortenUrl := recorder.Body.String()
	testData := []struct {
		shortUrl     string
		expectedCode int
		expectString string
	}{
		{shortenUrl, http.StatusOK, originalUrl},
		{"invalid-short-url", http.StatusUnprocessableEntity, ""},
		{originalUrl + "/invalid-short-url", http.StatusNotFound, ""},
	}

	for _, data := range testData {
		t.Run(data.shortUrl, func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/get?url="+data.shortUrl, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			assert.Equal(t, data.expectedCode, recorder.Code)
			assert.Equal(t, recorder.Body.String(), data.expectString)
		})

	}
}

func TestStats(t *testing.T) {
	router := SetUpRouter()
	originalUrl := "https://www.google.com"
	request, _ := http.NewRequest("POST", "/create?url="+originalUrl, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	shortenUrl := recorder.Body.String()
	testData := []struct {
		shortUrl     string
		expectedCode int
		expectString string
	}{
		{shortenUrl, http.StatusOK, originalUrl},
		{shortenUrl, http.StatusOK, originalUrl},
		{"invalid-short-url", http.StatusUnprocessableEntity, ""},
		{originalUrl + "/invalid-short-url", http.StatusNotFound, ""},
	}

	for index, data := range testData {
		testName := fmt.Sprintf("%s_%d", data.shortUrl, index)
		t.Run(testName, func(t *testing.T) {
			statsRequest, _ := http.NewRequest("GET", "/stats?url="+data.shortUrl, nil)
			statsRecorder := httptest.NewRecorder()
			router.ServeHTTP(statsRecorder, statsRequest)
			if data.expectedCode == http.StatusOK {
				jsonData := statsRecorder.Body.Bytes()
				var stats urlStorage.LinkStats
				err := json.Unmarshal(jsonData, &stats)
				assert.NoError(t, err)
				assert.Equal(t, index, stats.ClickCount)
				assert.Equal(t, data.expectString, stats.OriginalUrl)
				assert.Equal(t, data.shortUrl, stats.ShortUrl)
			}

			IncrementStats(data, router)
			assert.Equal(t, data.expectedCode, statsRecorder.Code)
		})

	}
}

func TestСreateCustomAlias(t *testing.T) {
	router := SetUpRouter()
	testData := []struct {
		requestUrl   string
		customAlias  string
		expectedCode int
		expectString bool
	}{
		{"https://www.google.com", "validAlias", http.StatusOK, true},
		{"https://www.example.com", "ivalid Alias", http.StatusUnprocessableEntity, false},
	}

	for _, data := range testData {
		t.Run(data.requestUrl, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/create?url="+data.requestUrl+"&custom_alias="+data.customAlias, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			assert.Equal(t, data.expectedCode, recorder.Code)
			assert.Equal(t, recorder.Body.String() != "", data.expectString)
		})

	}
}

func IncrementStats(data struct {
	shortUrl     string
	expectedCode int
	expectString string
}, router *gin.Engine) {
	getRequest, _ := http.NewRequest("GET", "/get?url="+data.shortUrl, nil)
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getRequest)
}
