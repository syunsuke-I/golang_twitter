# db
FROM postgres:14
ENV POSTGRES_PASSWORD postgres
RUN apt-get update && \
  apt-get clean && \
  rm -fr /var/lib/apt/lists/*

# builder
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# ソースの変更後にテストを実行するため
RUN go test ./models
RUN go build -o main ./main.go

# dev
FROM golang:1.21-alpine AS development
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]

