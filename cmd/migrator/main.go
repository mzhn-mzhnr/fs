package main

import (
	"errors"
	"mzhn/fileservice/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("cannot load .env")
	}
}

func main() {
	cfg := config.New()
	cs := cfg.PG.String()

	m, err := migrate.New(
		"file://migrations",
		cs,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}
}
