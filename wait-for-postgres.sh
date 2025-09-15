#!/bin/bash
set -e

host="$1"
shift
cmd="$@"

# Разбираем DATABASE_URL
# Формат: postgres://user:password@host:port/dbname?sslmode=disable
uri="$DATABASE_URL"

user=$(echo $uri | sed -E 's#postgres://([^:]+):.*#\1#')
password=$(echo $uri | sed -E 's#postgres://[^:]+:([^@]+)@.*#\1#')
db=$(echo $uri | sed -E 's#.*\/([^?]+).*#\1#')

export PGPASSWORD=$password

until psql -h "$host" -U "$user" -d "$db" -c '\q' 2>/dev/null; do
  echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - executing command"
exec $cmd
