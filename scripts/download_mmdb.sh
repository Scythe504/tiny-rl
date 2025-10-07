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

curl -L -u "${MAXMIND_ACCOUNT_ID}:${MAXMIND_LICENSE_KEY}" \
    "https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz" \
    -o /tmp/mmdb.tar.gz

tar -xzf /tmp/mmdb.tar.gz -C /tmp
find /tmp -name "*.mmdb" -exec mv {} /app/data/GeoLite2-Country.mmdb \;

rm -rf /tmp/mmdb.tar.gz /tmp/*

echo "GeoLite2-Country.mmdb downloaded successfully to /app/data/"
