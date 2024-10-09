FROM golang:1.23.2-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o my_app ./cmd/scheduler/main.go

FROM alpine:3.18

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/storage/scheduler.db

WORKDIR /app

COPY --from=builder /app/my_app /app/
COPY --from=builder /app/.env .
COPY --from=builder /app/web ./web 

RUN mkdir -p /app/storage

CMD ["./my_app"]