package urlStorage

import (
	"context"
	"fmt"
	"os"
	"strconv"

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
	query := `CREATE TABLE IF NOT EXISTS users (
		fullUrl TEXT  NOT NULL,
		shortUrl TEXT NOT NULL,
		clicks TEXT NOT NULL UNIQUE
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
	query := `INSRERT INTO users (fullUrl, shortUrl, clicks) VALUES ($1, $2, $3)`
	_, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, fullUrl, shortUrl, "0")
	return err
}

// Retrieves the original URL based on the shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetOriginalUrl(shortUrl string) (string, error) {
	query := `SELECT fullUrl from  users whhere shortUrl = $1`
	execResult, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, shortUrl)
	if err != nil {
		return "", err
	}

	return execResult.String(), nil
}

// Increments the click count for a given shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) IncrementClick(shortUrl string) error {
	query := `UPDATE users SET clicks = clicks + 1 WHERE shortUrl = $1`
	_, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, shortUrl)
	return err
}

// Retrieves the click count for a given shortened URL
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetClickCount(shortUrl string) (int, error) {
	query := `SELECT clicks from  users whhere shortUrl = $1`
	execResult, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, shortUrl)
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(execResult.String())
	return count, err
}

// Gets saved short URL for a given original URL, "" if not found
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) GetShortUrl(fillUrl string) (string, error) {
	query := `SELECT shortUrl from  users whhere fullUrl = $1`
	execResult, err := postgreSqlUrlStorage.connectionPool.Exec(context.Background(), query, fillUrl)
	if err != nil {
		return "", err
	}

	return execResult.String(), nil
}

// Closes the database connection pool
func (postgreSqlUrlStorage *PostgreSqlUrlStorage) Close() error {
	postgreSqlUrlStorage.connectionPool.Close()
	return nil
}
