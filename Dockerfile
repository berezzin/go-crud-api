FROM golang:1.24-rc-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crud-api .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/crud-api .
COPY --from=builder /app/cmd/goose ./cmd/goose
COPY --from=builder /app/db/migrations ./db/migrations

CMD ["/bin/sh", "-c", "cmd/goose postgres postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_HOST:$DB_PORT/$POSTGRES_DB -dir db/migrations up && ./crud-api"]
