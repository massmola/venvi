#!/bin/bash
set -e

DB_DIR=".pg_data"
LOG_FILE="postgres.log"

echo "🐘 setup local postgres..."

# 1. Initialize Data Directory
if [ ! -d "$DB_DIR" ]; then
    echo " -> Initializing database directory in $DB_DIR..."
    initdb -D "$DB_DIR" --no-locale --encoding=UTF8 > /dev/null
else
    echo " -> Database directory exists."
fi

# 2. Start PostgreSQL
if ! pg_ctl -D "$DB_DIR" status > /dev/null 2>&1; then
    echo " -> Starting PostgreSQL (Socket in /tmp)..."
    # -k specifies the socket directory
    pg_ctl -D "$DB_DIR" -l "$LOG_FILE" -o "-k /tmp" start
    sleep 3
else
    echo " -> PostgreSQL is already running."
fi

# 3. Create User and Database
# Use -h /tmp to connect via the unix socket we just configured
if ! psql -h /tmp template1 -tAc "SELECT 1 FROM pg_roles WHERE rolname='postgres'" | grep -q '1'; then
    echo " -> Creating superuser 'postgres'..."
    createuser -h /tmp -s postgres
    # Set password
    psql -h /tmp template1 -c "ALTER USER postgres WITH PASSWORD 'postgres';" > /dev/null
else
    echo " -> User 'postgres' exists."
fi

if ! psql -h /tmp template1 -tAc "SELECT 1 FROM pg_database WHERE datname='venvi'" | grep -q '1'; then
    echo " -> Creating database 'venvi'..."
    createdb -h /tmp -O postgres venvi
else
    echo " -> Database 'venvi' exists."
fi

echo "✅ Database ready!"
