deps:
	- go mod download
	- go mod tidy

run:
	- go run cmd/api/main.go

build:
	- go build -o guicpay cmd/api/main.go
