# Stage 1: Build
FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o taskvault .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/taskvault .

EXPOSE 8080 8946
CMD ["./taskvault"]