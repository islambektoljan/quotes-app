#!/bin/bash

# Скрипт для применения миграций вручную
echo "Applying database migrations..."

for migration_file in migrations/*.sql; do
    echo "Applying: $migration_file"
    PGPASSWORD=postgres psql -h localhost -U postgres -d quotes_db -f "$migration_file"
done

echo "All migrations applied successfully!"