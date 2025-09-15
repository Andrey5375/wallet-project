.PHONY: run test migrate

run:
	go run ./cmd/server

test:
	go test ./tests/...

migrate:
	psql $$DATABASE_URL -f migrations/001_init.sql
