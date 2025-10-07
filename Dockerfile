FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o build/tiny-rl cmd/api/main.go

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM debian:bookworm-slim

WORKDIR /app

# Install curl and bash, then clean up
RUN apt-get update && \
    apt-get install -y curl bash && \
    rm -rf /var/lib/apt/lists/*

# Copy binary from build directory
COPY --from=build /app/build ./build
COPY --from=build /app/scripts ./scripts
COPY --from=build /app/migrations ./migrations
# Copy goose binary from build stage
COPY --from=build /go/bin/goose /usr/local/bin/goose

RUN chmod +x /app/scripts/*.sh

CMD ["bash", "-c", "bash scripts/download_mmdb.sh && bash scripts/migrate.sh && ./build/tiny-rl"]