//go:build wireinject

package app

import (
	"context"
	"log/slog"
	"mzhn/fileservice/internal/config"
	"mzhn/fileservice/internal/providers/uuid"
	"mzhn/fileservice/internal/services/fileservice"
	"mzhn/fileservice/internal/storage/fs"
	"mzhn/fileservice/internal/storage/pg/filerepository"
	"mzhn/fileservice/internal/transport/http"
	"mzhn/fileservice/pkg/sl"
	"time"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		config.New,
		_context,
		_pgxpool,

		filerepository.New,
		uuid.New,
		fs.New,
		wire.Bind(new(fileservice.FileSaver), new(*filerepository.Repository)),
		wire.Bind(new(fileservice.FileProvider), new(*fs.FileStorage)),
		wire.Bind(new(fileservice.FileUploader), new(*fs.FileStorage)),
		wire.Bind(new(fileservice.IdProvider), new(*uuid.UuidProvider)),

		fileservice.New,

		_servers,
		newApp,
	))
}

func _servers(cfg *config.Config, fs *fileservice.FileService) []Server {
	servers := make([]Server, 0, 2)

	if cfg.Http.Enabled {
		servers = append(servers, http.New(cfg, fs))
	}

	return servers
}

func _context() context.Context {
	return context.Background()
}

func _pgxpool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
	fn := "_pgxpool"
	log := slog.With(slog.String("injector", fn))

	cs := cfg.PG.String()

	pool, err := pgxpool.New(ctx, cs)
	if err != nil {
		return nil, nil, err
	}

	start := time.Now()
	if err := pool.Ping(ctx); err != nil {
		return nil, nil, err
	}
	log.Info("connected to postgres", sl.Millis("took", time.Since(start)))

	return pool, func() {
		pool.Close()
	}, nil

}
