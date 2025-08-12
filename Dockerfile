# Dockerfile

# Stage 1: Build the Golang binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env ./.env

EXPOSE 8080

VOLUME /app/data
CMD ["./app"]