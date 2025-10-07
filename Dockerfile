FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with the correct path to main.go
RUN go build -o build/tiny-rl cmd/api/main.go

FROM debian:bookworm-slim

WORKDIR /app

# Install both curl and bash, then clean up
RUN apt-get update && \
    apt-get install -y curl bash && \
    rm -rf /var/lib/apt/lists/*

# Copy binary from build directory
COPY --from=build /app/build ./build
COPY --from=build /app/scripts ./scripts
COPY --from=build /app/migrations ./migrations

RUN chmod +x /app/scripts/*.sh

CMD ["bash", "-c", "bash scripts/download_mmdb.sh && bash scripts/migrate.sh && ./build/tiny-rl"]