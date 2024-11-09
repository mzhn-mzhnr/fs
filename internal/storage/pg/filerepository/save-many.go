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

func (r *Repository) SaveMany(ctx context.Context, records []*domain.FileRecord) error {
	fn := "pg.FileRepository.SaveMany"
	log := r.logger.With(sl.Method(fn))

	if len(records) == 0 {
		return nil
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer conn.Release()

	log.Debug("saving records", slog.Int("len", len(records)))

	builder := sq.
		Insert(pg.FilesTable).
		Columns("id", "name").
		PlaceholderFormat(sq.Dollar)

	for _, record := range records {
		builder = builder.Values(record.Id, record.Name)
	}

	query, args, err := builder.ToSql()
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
