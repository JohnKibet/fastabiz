# 1. Use a lightweight Go base
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# 2. Fetch Go module dependencies
COPY go.mod go.sum ./
RUN go mod download

# 3. Copy all source code
COPY . .

# 4. Build the binary
RUN go build -o fastabiz .

# 5. Final stage – smallest possible runtime image
FROM alpine:latest

WORKDIR /app

# 6. Copy binary from builder
COPY --from=builder /app/fastabiz .

# 7. Expose the port your app listens on
EXPOSE 8080

# 8. Run the binary
CMD ["./fastabiz"]
