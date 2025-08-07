package postgres

import (
    "context"
    "database/sql"

    _ "github.com/jackc/pgx/v5/stdlib"
    "link-shortener/internal/storage"
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS links (
    short_id VARCHAR(6) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
`

type PostgresRepository struct {
    db *sql.DB
}

func New(dsn string) (*PostgresRepository, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    if _, err := db.Exec(createTableQuery); err != nil {
        return nil, err
    }

    return &PostgresRepository{db: db}, nil
}

func (p *PostgresRepository) GetOriginalURL(shortID string) (string, bool) {
    var url string
    err := p.db.QueryRow("SELECT original_url FROM links WHERE short_id = $1", shortID).Scan(&url)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", false
        }
        return "", false
    }
    return url, true
}

func (p *PostgresRepository) AddShortURL(shortID, originalURL string) {
    _, _ = p.db.Exec(
        "INSERT INTO links (short_id, original_url) VALUES ($1, $2) ON CONFLICT (short_id) DO UPDATE SET original_url = $2",
        shortID, originalURL,
    )
}

func (p *PostgresRepository) Ping() error {
    return p.db.Ping()
}