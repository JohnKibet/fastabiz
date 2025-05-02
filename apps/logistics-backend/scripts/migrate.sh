#!/bin/bash

MIGRATIONS_DIR="./migrations"
DATABASE_URL="postgres://admin:secret@localhost:5432/logistics_db?sslmode=disable"

case "$1" in
  run)
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
    ;;
  down)
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down
    ;;
  drop)
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" drop -f
    ;;
  clean)
    echo "Cleaning all migration SQL files..."
    rm -f "$MIGRATIONS_DIR"/*.up.sql "$MIGRATIONS_DIR"/*.down.sql
    ;;
  drop_specific)
    shift
    for name in "$@"; do
      echo "Deleting migration: $name"
      rm -f "$MIGRATIONS_DIR"/*"$name"*.up.sql "$MIGRATIONS_DIR"/*"$name"*.down.sql
    done
    ;;
  create)
    shift
    for name in "$@"; do
      migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"
      sleep 1  # Ensure unique sequential numbering
    done
    ;;
  fix_dirty)
    echo "Forcing migration state to a clean version (default 0)"
    version=${2:-0}
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$version"
    ;;
  *)
    echo "Usage: ./scripts/migrate.sh [run|down|drop|clean|drop_specific name1 name2 ...|create name1 name2 ...|fix_dirty [version]]"
    ;;
esac
