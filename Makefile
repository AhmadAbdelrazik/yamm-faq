include .env

app-run:
	@go run ./cmd/api/main.go

psql:
	@psql -U postgres --dbname=$(DB_DATABASE)
