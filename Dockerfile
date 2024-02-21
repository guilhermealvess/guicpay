FROM golang:1.22

WORKDIR /app

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o guicpay cmd/api/main.go

EXPOSE 3000