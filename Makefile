deps:
	- go mod download
	- go mod tidy

run:
	- go run cmd/main.go

sql-generate:
	- go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	- sqlc generate