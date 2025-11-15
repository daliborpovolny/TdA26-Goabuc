# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

# Install tools
RUN apk add --no-cache bash git

WORKDIR /app

# Download dependencies first (cached layer)
COPY go.mod go.sum ./

COPY internal/database/db_file.db /db_file.db
COPY static /static

RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/tourbackend

# ---- Run Stage ----
FROM alpine:3.20

RUN apk add --no-cache bash netcat-openbsd

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy the folder of static files
COPY static ./static


# Environment variables
ENV PORT=3000
ENV PATH_TO_DB="db_file.db"
ENV STATIC_PATH="./static"

EXPOSE 3000

CMD ["./server"]
