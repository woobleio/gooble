#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  DELETE FROM package_creation;
  DELETE FROM package;
  DELETE FROM creation;
  DELETE FROM plan_user;
  DELETE FROM app_user;
EOSQL

pg_dump --file "/schema.sql" --username "wooble" --no-password --verbose --role "wooble" --format=p --encoding "UTF8" "wooble"
