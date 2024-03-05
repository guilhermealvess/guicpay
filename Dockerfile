FROM golang:1.22 AS builder

WORKDIR /app

COPY . .
RUN make build

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/guicpay /app/guicpay

COPY . .

RUN chmod +x /app/guicpay

RUN apk add --no-cache ca-certificates

RUN adduser -S gg -H && \
    chown -R gg: /app
USER gg

EXPOSE 3000