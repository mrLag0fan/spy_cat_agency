#!/bin/bash

# waiting for DB to be available
host="$1"
port="$2"
shift 2

# Додаємо кілька спроб підключитися
for i in {1..30}; do
    nc -z "$host" "$port" && break
    echo "Waiting for database connection..."
    sleep 2
done

if [[ $i -eq 30 ]]; then
    echo "Database is not available after 30 attempts, exiting"
    exit 1
fi

echo "Database is up and running, starting the application..."

# Запуск вашого додатка
exec "$@"
