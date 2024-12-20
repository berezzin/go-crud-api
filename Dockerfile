FROM golang:1.24-rc-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o crud-api .

CMD ["./crud-api"]
