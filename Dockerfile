# builder
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go

# dev
FROM golang:1.21-alpine AS development
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]