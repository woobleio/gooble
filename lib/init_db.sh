#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER user_test;
    CREATE DATABASE wooble_test;
    GRANT ALL PRIVILEGES ON DATABASE wooble_test TO user_test;
EOSQL
