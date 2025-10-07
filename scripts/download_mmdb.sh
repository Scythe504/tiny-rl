#!/bin/bash
set -e

echo "Downloading GeoLite2-Country.mmdb ..."

mkdir -p /app/data

if [ -z "$MAXMIND_ACCOUNT_ID" ] || [ -z "$MAXMIND_LICENSE_KEY" ]; then
    echo "Missing MaxMind credentials!"
    echo "Please provide both environment variables:"
    echo "  MAXMIND_ACCOUNT_ID=<your-account-id>"
    echo "  MAXMIND_LICENSE_KEY=<your-license-key>"
    exit 1
fi

echo "Downloading from MaxMind..."

# Create a temporary directory for this operation
TEMP_DIR=$(mktemp -d)
trap 'rm -rf $TEMP_DIR' EXIT

curl -L -u "${MAXMIND_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY}" \
    "https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz" \
    -o "$TEMP_DIR/mmdb.tar.gz"

echo "Extracting archive..."
tar -xzf "$TEMP_DIR/mmdb.tar.gz" -C "$TEMP_DIR"

# Find and move the .mmdb file
MMDB_FILE=$(find "$TEMP_DIR" -name "*.mmdb" -type f | head -n 1)

if [ -z "$MMDB_FILE" ]; then
    echo "Error: Could not find .mmdb file in downloaded archive"
    exit 1
fi

mv "$MMDB_FILE" /app/data/GeoLite2-Country.mmdb

echo "GeoLite2-Country.mmdb downloaded successfully to /app/data/"