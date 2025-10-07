#!/bin/bash

echo "Running migrations..."
cd /app/migrations || return 1

goose "$POSTGRES_CONN_URL" up
echo "Database migration complete"
