package database

import (
	"database/sql"
	"fmt"

	"github.com/Ruseigha/sendlink/internal/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
    db *sql.DB
}

func NewPostgresDB(connectionString string) (*PostgresDB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) InitSchema() error {
    schema := `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        original_url TEXT NOT NULL,
        short_code VARCHAR(10) UNIQUE NOT NULL,
        clicks INTEGER DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
    `
    _, err := p.db.Exec(schema)
    return err
}

func (p *PostgresDB) CreateURL(url *models.URL) error {
    query := `
        INSERT INTO urls (original_url, short_code, expires_at)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
    return p.db.QueryRow(query, url.OriginalURL, url.ShortCode, url.ExpiresAt).
        Scan(&url.ID, &url.CreatedAt)
}

func (p *PostgresDB) GetURLByShortCode(shortCode string) (*models.URL, error) {
    url := &models.URL{}
    query := `
        SELECT id, original_url, short_code, clicks, created_at, expires_at
        FROM urls
        WHERE short_code = $1
    `
    err := p.db.QueryRow(query, shortCode).Scan(
        &url.ID, &url.OriginalURL, &url.ShortCode,
        &url.Clicks, &url.CreatedAt, &url.ExpiresAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("URL not found")
    }
    return url, err
}

func (p *PostgresDB) IncrementClicks(shortCode string) error {
    query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1`
    _, err := p.db.Exec(query, shortCode)
    return err
}

func (p *PostgresDB) GetStats() (*models.Stats, error) {
    stats := &models.Stats{}
    query := `SELECT COUNT(*), COALESCE(SUM(clicks), 0) FROM urls`
    err := p.db.QueryRow(query).Scan(&stats.TotalURLs, &stats.TotalClicks)
    return stats, err
}

func (p *PostgresDB) Close() error {
    return p.db.Close()
}