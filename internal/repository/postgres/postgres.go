package postgres

import (
	"database/sql"
	"time"
    "context"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
й
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return nil, err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (p *PostgresRepository) GetOriginalURL(shortID string) (string, bool) {
	var url string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.db.QueryRowContext(ctx, "SELECT original_url FROM links WHERE short_id = $1", shortID).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		return "", false
	}
	return url, true
}

func (p *PostgresRepository) AddShortURL(shortID, originalURL string) error {
    _, err := p.db.Exec(
        "INSERT INTO links (short_id, original_url) VALUES ($1, $2) ON CONFLICT (short_id) DO UPDATE SET original_url = $2",
        shortID, originalURL,
    )
    return err
}

func (p *PostgresRepository) Ping() error {
    return p.db.Ping()
}