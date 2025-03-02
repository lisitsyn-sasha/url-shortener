package storage

import "context"

type URLSaver interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error)
}

type URLDeleter interface {
	DeleteUrl(ctx context.Context, alias string) (int64, error)
}
