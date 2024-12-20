FROM golang:1.24-rc-alpine

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crud-api .

CMD ["/bin/sh", "-c", "goose postgres postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_HOST:$DB_PORT/$POSTGRES_DB -dir db/migrations up && ./crud-api"]
