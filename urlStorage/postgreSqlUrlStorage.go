package urlStorage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSqlUrlStorage is a PostgreSQL implementation of the UrlStorage interface
type PostgreSqlUrlStorage struct {
	connectionPool *pgxpool.Pool
}

func NewPostgreSqlUrlStorage() (PostgreSqlUrlStorage, error) {
	postgreSqlUrlStorage := PostgreSqlUrlStorage{}
	dsn := os.Getenv("POSTGRES_DATA_SOURCE_NAME")
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return postgreSqlUrlStorage, err
	}
	postgreSqlUrlStorage.connectionPool = dbPool
	query := `CREATE TABLE IF NOT EXISTS urls (
		fullUrl TEXT  NOT NULL UNIQUE,
		shortUrl TEXT NOT NULL UNIQUE,
		clicks INTEGER NOT NULL DEFAULT 0,
		lastAccessedAt TIMESTAMPTZ,
		createdAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = dbPool.Exec(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		return postgreSqlUrlStorage, err
	}

	return postgreSqlUrlStorage, err
}

// Saves the mapping between the original URL and the shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) Store(ctx context.Context, fullUrl string, shortUrl string) error {
	query := `INSERT INTO urls (fullUrl, shortUrl, clicks) 
	VALUES ($1, $2, $3) 
	ON CONFLICT DO NOTHING`
	_, err := postgreSqlUrlStorage.connectionPool.Exec(ctx, query, fullUrl, shortUrl, "0")
	return err
}

// Retrieves the original URL based on the shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	query := `SELECT fullUrl from  urls where shortUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(ctx, query, shortUrl)
	var fullUrl string
	err := execResult.Scan(&fullUrl)
	if err != nil {
		return "", err
	}

	return fullUrl, nil
}

// Retrieves the click count for a given shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetStats(ctx context.Context, shortUrl string) (LinkStats, error) {
	query := `SELECT clicks, fullUrl, lastAccessedAt, createdAt from  urls where shortUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(ctx, query, shortUrl)
	var stats LinkStats
	err := execResult.Scan(&stats.ClickCount, &stats.OriginalUrl, &stats.LastAccesedAt, &stats.CreatedAt)
	if err != nil {
		return LinkStats{}, err
	}

	return stats, err
}

// Gets saved short URL for a given original URL, "" if not found
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetShortUrl(ctx context.Context, fillUrl string) (string, error) {
	query := `SELECT shortUrl from  urls where fullUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(ctx, query, fillUrl)
	var shortUrl string
	err := execResult.Scan(&shortUrl)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

// Deletes links that are older than the specified duration
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) DeleteLinksOlderThan(ctx context.Context, olderThan time.Duration) {
	oldTime := time.Now().Add(-olderThan)
	query := `DELETE from urls where createdAt < $1`
	postgreSqlUrlStorage.connectionPool.Exec(ctx, query, oldTime)
}

// Closes the database connection pool
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) Close() error {
	postgreSqlUrlStorage.connectionPool.Close()
	return nil
}

// Update stats for short URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) UpdateStats(ctx context.Context, shortUrl string) error {
	query := `UPDATE urls 
				SET 
				clicks = clicks + 1,
				lastAccessedAt = CURRENT_TIMESTAMP
				WHERE shortUrl = $1`
	_, err := postgreSqlUrlStorage.connectionPool.Exec(ctx, query, shortUrl)
	return err
}
