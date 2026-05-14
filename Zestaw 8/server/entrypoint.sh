#!/bin/sh
set -e

# Seed only when database file is missing or empty.
if [ ! -s /app/shop.db ]; then
  echo "Seeding database..."
  /app/seed
fi

exec /app/server