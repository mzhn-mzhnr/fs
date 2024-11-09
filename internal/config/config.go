package config

import (
	"fmt"
	"log/slog"
	"mzhn/fileservice/pkg/prettyslog"
	"mzhn/fileservice/pkg/sl"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Name    string `env:"APP_NAME" env-default:"mzhn-fileservice"`
	Version string `env:"APP_VERSION" env-default:"local"`
	Env     string `env:"APP_ENV" env-default:"local"`
}

type Http struct {
	Enabled bool   `env:"HTTP_ENABLED" env-default:"true"`
	Host    string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port    int    `env:"HTTP_PORT" env-default:"8080"`
}

type Postgres struct {
	Host string `env:"DATABASE_HOST" env-default:"localhost"`
	Port int    `env:"DATABASE_PORT" env-default:"5432"`
	User string `env:"DATABASE_USER" env-default:"postgres"`
	Pass string `env:"DATABASE_PASS" env-default:"postgres"`
	Name string `env:"DATABASE_NAME" env-default:"files"`
	SSL  string `env:"DATABASE_SSL" env-default:"disable"`
}

func (p *Postgres) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", p.User, p.Pass, p.Host, p.Port, p.Name, p.SSL)
}

type FileStorage struct {
	Path string `env:"FS_PATH" env-default:"./static/"`
}

type Config struct {
	App  App
	Http Http
	FS   FileStorage
	PG   Postgres
}

func New() *Config {
	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		slog.Error("error when reading env", sl.Err(err))
		header := fmt.Sprintf("%s - %s", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION"))

		usage := cleanenv.FUsage(os.Stdout, config, &header)
		usage()

		os.Exit(-1)
	}

	setupLogger(config)

	return config
}

func setupLogger(cfg *Config) {
	var log *slog.Logger

	switch cfg.App.Env {
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(prettyslog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	slog.SetDefault(log)
}
