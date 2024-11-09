package filerepository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/internal/storage"
	"mzhn/fileservice/internal/storage/pg"
	"mzhn/fileservice/pkg/sl"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Save(ctx context.Context, record domain.FileRecord) error {
	fn := "pg.FileRepository.Save"
	log := r.logger.With(sl.Method(fn), slog.String("id", record.Id), slog.String("filename", record.Name))

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer conn.Release()

	query, args, err := sq.
		Insert(pg.FilesTable).
		Columns("id", "name").
		Values(record.Id, record.Name).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	log = log.With(slog.String("query", query), slog.Any("args", args))

	log.Debug("executing")
	start := time.Now()
	if _, err := conn.Exec(ctx, query, args...); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) {
			switch e.Code {
			case "23505":
				return fmt.Errorf("%s: %w", fn, storage.ErrFileDuplicate)
			case "22P02":
				return fmt.Errorf("%s: %w", fn, storage.ErrFileBadId)
			}
		}
		log.Error("failed to execute query")
		return fmt.Errorf("%s: %w", fn, err)
	}
	log.Debug("executed", sl.Millis("took", time.Since(start)))

	return nil
}
