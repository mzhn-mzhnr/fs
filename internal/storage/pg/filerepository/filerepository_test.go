package filerepository_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mzhn/fileservice/internal/config"
	"mzhn/fileservice/internal/domain"
	"mzhn/fileservice/internal/providers/uuid"
	"mzhn/fileservice/internal/storage"
	def "mzhn/fileservice/internal/storage/pg/filerepository"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	dbuser = "postgres"
	dbpass = "postgres"
	dbhost = "localhost"
	dbport = 5432
	dbname = "files_test"
	defdb  = "postgres"
)

func conf(t *testing.T, name string) *config.Config {
	t.Helper()
	return &config.Config{
		PG: config.Postgres{
			Host: dbhost,
			Port: dbport,
			User: dbuser,
			Pass: dbpass,
			Name: name,
			SSL:  "disable",
		},
	}
}

func connect(t *testing.T, ctx context.Context, name string) (*pgxpool.Pool, error) {
	t.Helper()

	cfg := conf(t, name)

	cs := cfg.PG.String()
	m, err := migrate.New(
		"file://../../../../migrations",
		cs,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			require.NoError(t, err)
		}
	}

	pool, err := pgxpool.New(ctx, cfg.PG.String())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func up(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	pool, err := connect(t, ctx, dbname)
	if err == nil {
		return pool
	}

	cfg := conf(t, defdb)
	pool, err = pgxpool.New(ctx, cfg.PG.String())
	require.NoError(t, err)

	if err := pool.Ping(ctx); err != nil {
		require.NoError(t, err)
	}

	if _, err := pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\"", dbname)); err != nil {
		require.NoError(t, err)
	}
	pool.Close()

	pool, _ = connect(t, ctx, dbname)
	return pool
}

func suite(t *testing.T) (context.Context, *def.Repository) {
	t.Helper()

	ctx := context.Background()

	pool := up(t, ctx)

	slog.SetLogLoggerLevel(slog.LevelDebug)

	repo := def.New(pool)
	return ctx, repo
}

func id() string {
	provider := uuid.New()
	return provider.Provide()
}

func TestSave(t *testing.T) {
	t.Parallel()
	ctx, repo := suite(t)

	cases := []struct {
		name     string
		id       string
		filename string
	}{
		{"txt", id(), "lorem ipsum.txt"},
		{"pdf", id(), "test.pdf"},
		{"png", id(), "test.png"},
		{"cyrillic", id(), "Отчёт.docx"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := repo.Save(ctx, &domain.FileRecord{Id: c.id, Name: c.filename})
			require.NoError(t, err)
		})
	}
}

func TestSaveDuplicate(t *testing.T) {
	t.Parallel()
	ctx, repo := suite(t)

	id := id()
	err := repo.Save(ctx, &domain.FileRecord{Id: id, Name: "data"})
	require.NoError(t, err)

	err = repo.Save(ctx, &domain.FileRecord{Id: id, Name: "data2"})
	require.ErrorIs(t, err, storage.ErrFileDuplicate)
}

func TestSaveNotUuid(t *testing.T) {
	t.Parallel()

	ctx, repo := suite(t)
	err := repo.Save(ctx, &domain.FileRecord{Id: "not-uuid", Name: "data"})
	require.ErrorIs(t, err, storage.ErrFileBadId)
}

func TestSaveMany(t *testing.T) {
	t.Parallel()

	ctx, repo := suite(t)

	n := 30
	records := make([]*domain.FileRecord, 0, n)
	for i := range n {
		id := id()
		records = append(records, &domain.FileRecord{Id: id, Name: fmt.Sprintf("data-%d", i)})
	}

	err := repo.SaveMany(ctx, records)
	require.NoError(t, err)
}
