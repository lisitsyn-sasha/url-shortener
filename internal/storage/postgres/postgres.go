package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"url-shortener/constants"
)

func New(connString string) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.PostgresNew, err)
	}

	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", constants.PostgresNew, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func createTables(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS url(
            id SERIAL PRIMARY KEY,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("%s: %w", constants.PostgresCreateTable, err)
	}

	_, err = db.Exec(context.Background(), `
        CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
    `)
	if err != nil {
		return fmt.Errorf("%s: %w", constants.PostgresCreateTable, err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	if alias == "" {
		return "", fmt.Errorf("%s: alias must not be empty", constants.PostgresGetUrl)
	}

	var url string

	query := `SELECT url FROM url WHERE alias = $1`

	err := s.db.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: alias not found", constants.PostgresGetUrl)
		}
		return "", fmt.Errorf("%s: %w", constants.PostgresGetUrl, err)
	}

	return url, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	if urlToSave == "" || alias == "" {
		return 0, fmt.Errorf("%s: url and alias must not be empty", constants.PostgresSaveUrl)
	}

	var id int64
	query := `
		INSERT INTO url(url, alias) VALUES($1, $2)
		ON CONFLICT (alias) DO UPDATE SET id = url.id
		RETURNING id
	`

	err := s.db.QueryRow(ctx, query, urlToSave, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", constants.PostgresSaveUrl, err)
	}

	return id, nil
}

func (s *Storage) DeleteUrl(ctx context.Context, alias string) (int64, error) {
	if alias == "" {
		return 0, fmt.Errorf("%: alias must not be empty", constants.PostgresDeleteUrl)
	}

	var id int64
	query := `
    	DELETE FROM url WHERE alias = $1 RETURNING id
	`

	err := s.db.QueryRow(ctx, query, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", constants.PostgresDeleteUrl, err)
	}

	return id, nil
}
