package urlStorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func GetUrlStoragesToTest(t *testing.T) []UrlStorage {
	inMemoryStorage := NewInMemoryUrlStorage()
	dsn := SetupPostgresContainer(t)
	postgresStorage, _ := NewPostgreSqlUrlStorage(dsn)
	urlstorages := []UrlStorage{
		&inMemoryStorage,
		&postgresStorage,
	}

	return urlstorages
}

func SetupPostgresContainer(t *testing.T) string {
	pgContainer, err := postgres.Run(t.Context(),
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies())
	if err != nil {
		panic("Unable to start PostgreSQL container")
	}
	t.Cleanup(func() {
		pgContainer.Terminate(t.Context())
	})

	connString, err := pgContainer.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		panic("Unable to get connection string for PostgreSQL container")
	}
	return connString
}

func TestUrlStorages(t *testing.T) {
	urlStoragesToTest := GetUrlStoragesToTest(t)
	for _, urlStorageToTest := range urlStoragesToTest {
		testedType := fmt.Sprintf("%T", urlStorageToTest)
		t.Run(testedType, func(t *testing.T) {
			testStore(t, &urlStorageToTest)
			testGetOriginalUrl(t, &urlStorageToTest)
			testGetStats(t, &urlStorageToTest)
			testGetShortUrl(t, &urlStorageToTest)
			testDeleteLinksOlderThan(t, &urlStorageToTest)
		})
	}
}

func testStore(t *testing.T, urlStorage *UrlStorage) {
	urlToShorten := "https://www.example.com"
	shortUrl := "abc123"
	err := (*urlStorage).Store(t.Context(), urlToShorten, shortUrl)
	assert.Nil(t, err)
}

func testGetOriginalUrl(t *testing.T, urlStorage *UrlStorage) {
	urlToShorten := "https://www.example.com"
	shortUrl := "abc123"
	(*urlStorage).Store(t.Context(), urlToShorten, shortUrl)
	testData := []struct {
		shortUrl    string
		expectedUrl string
		expectError bool
	}{
		{shortUrl: shortUrl, expectedUrl: urlToShorten, expectError: false},
		{shortUrl: "nonExistentShortUrl", expectedUrl: "", expectError: true},
	}

	for _, testCase := range testData {
		t.Run(testCase.shortUrl, func(t *testing.T) {
			originalUrl, err := (*urlStorage).GetOriginalUrl(t.Context(), testCase.shortUrl)
			assert.Equal(t, testCase.expectedUrl, originalUrl)
			assert.Equal(t, testCase.expectError, err != nil)
		})
	}
}

func testGetStats(t *testing.T, urlStorage *UrlStorage) {
	urlToShorten := "https://www.example.com"
	shortUrl := "abc123"
	(*urlStorage).Store(t.Context(), urlToShorten, shortUrl)
	t.Run("Existing url, no clicks", func(t *testing.T) {
		stats, err := (*urlStorage).GetStats(t.Context(), shortUrl)
		assert.Nil(t, err)
		assert.Equal(t, shortUrl, stats.ShortUrl)
		assert.Equal(t, urlToShorten, stats.OriginalUrl)
		assert.Equal(t, 0, stats.ClickCount)
	})
	t.Run("Existing url, 2 clicks", func(t *testing.T) {
		(*urlStorage).UpdateStats(t.Context(), shortUrl)
		(*urlStorage).UpdateStats(t.Context(), shortUrl)
		stats, err := (*urlStorage).GetStats(t.Context(), shortUrl)
		assert.Nil(t, err)
		assert.Equal(t, shortUrl, stats.ShortUrl)
		assert.Equal(t, urlToShorten, stats.OriginalUrl)
		assert.Equal(t, 2, stats.ClickCount)
		assert.NotNil(t, stats.CreatedAt)
		assert.NotNil(t, stats.LastAccesedAt)
	})

	t.Run("Not Existing url", func(t *testing.T) {
		nonExistingShortUrl := "nonExistentShortUrl"
		_, err := (*urlStorage).GetStats(t.Context(), nonExistingShortUrl)
		assert.NotNil(t, err)
	})
}

func testGetShortUrl(t *testing.T, urlStorage *UrlStorage) {
	originalUrl := "https://www.example.com"
	shortUrl := "abc123"
	(*urlStorage).Store(t.Context(), originalUrl, shortUrl)
	testData := []struct {
		shortUrl    string
		originalUrl string
		expectError bool
	}{
		{shortUrl: shortUrl, originalUrl: originalUrl, expectError: false},
		{shortUrl: "", originalUrl: "nonExistentShortUrl", expectError: true},
	}

	for _, testCase := range testData {
		t.Run(testCase.shortUrl, func(t *testing.T) {
			actualShortUrl, err := (*urlStorage).GetShortUrl(t.Context(), testCase.originalUrl)
			assert.Equal(t, testCase.shortUrl, actualShortUrl)
			assert.Equal(t, testCase.expectError, err != nil)
		})
	}
}

func testDeleteLinksOlderThan(t *testing.T, urlStorage *UrlStorage) {
	originalUrl := "https://www.example.com"
	shortUrl := "abc123"
	(*urlStorage).Store(t.Context(), originalUrl, shortUrl)
	t.Run("Delete old links - no links to delete", func(t *testing.T) {
		(*urlStorage).DeleteLinksOlderThan(t.Context(), time.Hour*2)
		actualOriginalUrl, err := (*urlStorage).GetOriginalUrl(t.Context(), shortUrl)
		assert.Equal(t, originalUrl, actualOriginalUrl)
		assert.Nil(t, err)
	})
	t.Run("Delete old links - link should be deleted", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond)
		(*urlStorage).DeleteLinksOlderThan(t.Context(), time.Millisecond)
		actualOriginalUrl, err := (*urlStorage).GetOriginalUrl(t.Context(), shortUrl)
		assert.Equal(t, "", actualOriginalUrl)
		assert.NotNil(t, err)
	})

}
