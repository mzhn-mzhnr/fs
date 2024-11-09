package filerepository

import (
	"log/slog"
	"mzhn/fileservice/pkg/sl"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:   pool,
		logger: slog.With(sl.Module("pg.FileRepository")),
	}
}
