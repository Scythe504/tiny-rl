#!/bin/bash
set -e

echo "Running migrations..."

if [ -z "$POSTGRES_CONN_URL" ]; then
    echo "Error: POSTGRES_CONN_URL environment variable is not set"
    exit 1
fi

if ! command -v goose &> /dev/null; then
    echo "Error: goose is not installed"
    exit 1
fi

cd /app/migrations || exit 1

goose postgres "$POSTGRES_CONN_URL" up

echo "Database migration complete"