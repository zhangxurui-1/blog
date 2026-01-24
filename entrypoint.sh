#!/bin/sh
set -e

echo "Initializing/updating database tables..."
/app/server --init-db-table

echo "Initializing Elasticsearch index..."
/app/server --init-es-index-default

echo "Starting server..."
exec /app/server
