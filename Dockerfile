FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM golang:1.24.2

WORKDIR /app

COPY --from=builder /app/main .

COPY .env ./

EXPOSE 8080

CMD ["./main"]