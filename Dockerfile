FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

RUN apk update && \
    apk add --no-cache tzdata && \
    apk add nano

WORKDIR /app

COPY --from=builder /app/guicpay /app/.

COPY docs/ /app/docs

RUN adduser -S guicpay -H && \
    chown -R guicpay: /app
USER guicpay

EXPOSE 3000
EXPOSE 5000

CMD [ "./guicpay" ]