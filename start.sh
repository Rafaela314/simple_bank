#!/bin/sh

set -e

echo "run db migrations"
source app.env
/app/migrate -path /app/migration -database "$DB_Source" -verbose up

echo "start the app"
exec "$@"