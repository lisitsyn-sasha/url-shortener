package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func createTables(db *pgxpool.Pool) error {
	const op = "storage.postgres.createTables"

	_, err := db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS url(
            id SERIAL PRIMARY KEY,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(context.Background(), `
        CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
    `)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	if urlToSave == "" || alias == "" {
		return 0, fmt.Errorf("%s: url and alias must not be empty", op)
	}

	var id int64
	query := `
		INSERT INTO url(url, alias) VALUES($1, $2)
		ON CONFLICT (alias) DO UPDATE SET id = url.id
		RETURNING id
	`

	err := s.db.QueryRow(ctx, query, urlToSave, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	if alias == "" {
		return "", fmt.Errorf("%s: alias must not be empty", op)
	}

	var url string

	query := `SELECT url FROM url WHERE alias = $1`

	err := s.db.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: alias not found", op)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// TODO: func (s *Storage) DeleteURL(alias string) error {}
