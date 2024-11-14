#!/bin/sh

set -e

echo "Run db migrations"
#config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName
/app/migrate -path /app/migrations -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -verbose up

echo "Starting the application"
exec "$@"