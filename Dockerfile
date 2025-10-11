FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from the root
COPY go.mod go.sum ./
RUN go mod download

# Copy all backend
COPY ./backend ./backend

# Compile the application
RUN go build -o main ./backend

FROM alpine:latest
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/main .

# Copy the frontend files to be served by the backend
COPY ./frontend ./frontend

EXPOSE 8080
CMD ["./main"]

# docker compose up --build
# docker compose ps
#  stop conteiner = docker compose down
# docker compose up -d --remove-orphans --> up new projct and to clean the containers that nnot use more
