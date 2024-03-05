LINUX_AMD64 = CGO_ENABLED=0 GOOS=linux GOARCH=amd64

deps:
	- go mod download
	- go mod tidy

run:
	- go run cmd/api/main.go

build:
	$(LINUX_AMD64) go build -o guicpay ./cmd/api/main.go

docker-run:
	- docker-compose up -d

ping:
	- curl http://localhost:8080/ping

gen-proto:
	- protoc --go_out=. --go-grpc_out=. pkg/pb/*.proto

test:
	go test -covermode=count -coverprofile=count.out ./...