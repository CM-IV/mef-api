#!/bin/sh

set -e

# echo "running migrations..."
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "running migrations..."
migrate -path db/migration -database "$DB_SOURCE" -verbose up

echo "Run air config..."
air -c .air/.air.toml

echo "start the application..."
exec "$@"
