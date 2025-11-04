FROM golang:1.23-alpine AS builder

WORKDIR /app/backend

# Copy go mod and sum files from backend directory
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source files
COPY backend/ .

# Compile the application
RUN go build -o /main .

FROM alpine:latest
WORKDIR /app

# Copy the binary 
COPY --from=builder /main /app/main

# Copy frontend files
COPY ./frontend ./frontend

EXPOSE 8080
CMD ["./main"]
