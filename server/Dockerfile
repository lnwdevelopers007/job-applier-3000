# 1. Go language set up
FROM golang:1.24 AS builder
WORKDIR /app

# Copy go.mod, go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy everything else
COPY . .

# Go build
RUN CGO_ENABLED=0 go build -o server .

# 2. Run stage
FROM alpine:3.19
WORKDIR /app

# Copy binary
COPY --from=builder /app/server /app/server

# Copy env file
COPY .env.production .env

CMD ["./server"]
