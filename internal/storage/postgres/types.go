package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	db *pgxpool.Pool
}
