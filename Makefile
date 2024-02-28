deps:
	- go mod download
	- go mod tidy

run:
	- go run cmd/api/main.go

build:
	- go build -o guicpay cmd/api/main.go

docker-run:
	- docker-compose up -d

ping:
	- curl http://localhost:8080/api/ping

gen-proto:
	- protoc --go_out=. --go-grpc_out=. pkg/pb/*.proto