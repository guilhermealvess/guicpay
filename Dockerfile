FROM golang:1.22 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o guicpay cmd/api/main.go

FROM alpine:latest

WORKDIR /app
RUN apk --update add sqlite && \
    rm -rf /var/cache/apk/*

COPY --from=builder /app/guicpay .
COPY guicpay.db .

RUN adduser -S app -H && \
    chown -R app: /app
USER app

EXPOSE 3000
