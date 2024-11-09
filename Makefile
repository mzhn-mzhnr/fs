ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen lint
	go build -o ./bin/app.exe ./cmd/app

run: build
	./bin/app -local

wire-gen:
	go generate ./internal/app

gen: wire-gen

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out

lint:
	golangci-lint run

# migrate.up:
#         migrate -path ./migrations -database "postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" up

# migrate.down:
#         migrate -path ./migrations -database "postgres://$(DATABASE_USER):$(DATABASE_PASS)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" down
