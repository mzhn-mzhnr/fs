package main

import (
	"context"
	"flag"
	"fmt"

	"mzhn/fileservice/internal/app"

	"github.com/joho/godotenv"
)

var (
	local bool
)

func init() {
	flag.BoolVar(&local, "local", false, "run locally")

}

func main() {
	flag.Parse()

	if local {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	instance, cleanup, err := app.New()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	instance.Run(context.Background())
}
