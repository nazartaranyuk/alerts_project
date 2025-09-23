#!/usr/bin/env bash
set -euo pipefail

PG_CONT=${PG_CONT:-pg}
PG_USER=${PG_USER:-myuser}
PG_PASS=${PG_PASS:-secret}
PG_HOST=${PG_HOST:-localhost}
PG_PORT=${PG_PORT:-5432}
APP_DB=${APP_DB:-alarms}

docker exec -e PGPASSWORD="$PG_PASS" -i "$PG_CONT" \
  psql -U "$PG_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${APP_DB}'" | grep -q 1 \
  || docker exec -e PGPASSWORD="$PG_PASS" -i "$PG_CONT" \
       psql -U "$PG_USER" -d postgres -c "CREATE DATABASE ${APP_DB}"

docker exec -i "$PG_CONT" bash -lc "cat >/tmp/database_scheme.sql" < database_scheme.sql
docker exec -e PGPASSWORD="$PG_PASS" -i "$PG_CONT" \
  psql -U "$PG_USER" -d "$APP_DB" -f /tmp/database_scheme.sql

echo "DB '${APP_DB}' is ready."
