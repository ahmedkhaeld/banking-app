#!/bin/sh

set -e

echo "Waiting for database to be ready..."
/app/wait-for-it.sh db:5432 --

echo "Starting the app..."
exec /app/main