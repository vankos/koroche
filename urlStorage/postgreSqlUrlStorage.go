package urlStorage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSqlUrlStorage is a PostgreSQL implementation of the UrlStorage interface
type PostgreSqlUrlStorage struct {
	connectionPool *pgxpool.Pool
}

func NewPostgreSqlUrlStorage() (PostgreSqlUrlStorage, error) {
	postgreSqlUrlStorage := PostgreSqlUrlStorage{}
	dsn := "postgres://postgres:secret@localhost:5432/postgres"
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return postgreSqlUrlStorage, err
	}
	postgreSqlUrlStorage.connectionPool = dbPool
	query := `CREATE TABLE IF NOT EXISTS urls (
		fullUrl TEXT  NOT NULL,
		shortUrl TEXT NOT NULL,
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
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) Store(fullUrl string, shortUrl string) error {
	query := `INSERT INTO urls (fullUrl, shortUrl, clicks) VALUES ($1, $2, $3)`
	_, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, fullUrl, shortUrl, "0")
	return err
}

// Retrieves the original URL based on the shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetOriginalUrl(shortUrl string) (string, error) {
	query := `SELECT fullUrl from  urls where shortUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(context.Background(), query, shortUrl)
	var fullUrl string
	err := execResult.Scan(&fullUrl)
	if err != nil {
		return "", err
	}

	UpdateStats(shortUrl, postgreSqlUrlStorage.connectionPool)
	return fullUrl, nil
}

// Increments the click count for a given shortened URL
func UpdateStats(shortUrl string, connectionPool *pgxpool.Pool) error {
	query := `UPDATE urls 
				SET 
				clicks = clicks + 1,
				lastAccessedAt = CURRENT_TIMESTAMP
				WHERE shortUrl = $1`
	_, err := connectionPool.Exec(context.Background(), query, shortUrl)
	return err
}

// Retrieves the click count for a given shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetStats(shortUrl string) (LinkStats, error) {
	query := `SELECT clicks, fullUrl, lastAccessedAt, createdAt from  urls where shortUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(context.Background(), query, shortUrl)
	var stats LinkStats
	err := execResult.Scan(&stats.ClickCount, &stats.OriginalUrl, &stats.LasAccesedAt, &stats.CreatedAt)
	if err != nil {
		return LinkStats{}, err
	}

	return stats, err
}

// Gets saved short URL for a given original URL, "" if not found
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetShortUrl(fillUrl string) (string, error) {
	query := `SELECT shortUrl from  urls where fullUrl = $1`
	execResult := postgreSqlUrlStorage.connectionPool.QueryRow(context.Background(), query, fillUrl)
	var shortUrl string
	err := execResult.Scan(&shortUrl)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

// Closes the database connection pool
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) Close() error {
	postgreSqlUrlStorage.connectionPool.Close()
	return nil
}
