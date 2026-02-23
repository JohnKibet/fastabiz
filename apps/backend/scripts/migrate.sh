#!/bin/bash
set -e

MIGRATIONS_DIR="./migrations"
DATABASE_URL="postgres://admin:secret@localhost:5432/fastabiz_db?sslmode=disable"

# --- Ensure migrate is installed ---
ensure_migrate() {
  if ! command -v migrate >/dev/null 2>&1; then
    echo "⚠️  'migrate' not found. Installing..."

    if ! command -v go >/dev/null 2>&1; then
      echo "❌ Go is not installed. Please install Go first."
      exit 1
    fi

    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

    # Add GOPATH/bin to PATH if needed
    export PATH="$PATH:$(go env GOPATH)/bin"

    if ! command -v migrate >/dev/null 2>&1; then
      echo "❌ Failed to install migrate. Ensure \$GOPATH/bin is in your PATH."
      exit 1
    fi

    echo "✅ migrate installed successfully."
  fi
}

case "$1" in
  run)
    ensure_migrate
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
    ;;

  run_table)
    shift
    for name in "$@"; do
      echo "Running migrations for table: $name"
      for file in $(ls "$MIGRATIONS_DIR" | grep "$name" | grep "\.up\.sql"); do
        echo "Applying $file"
        psql "$DATABASE_URL" -f "$MIGRATIONS_DIR/$file"
      done
    done
    ;;

  down)
    ensure_migrate
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down
    ;;

  drop)
    ensure_migrate
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" drop -f
    ;;

  drop_table)
    shift
    for table in "$@"; do
      echo "Dropping table: $table"
      psql "$DATABASE_URL" -c "DROP TABLE IF EXISTS $table CASCADE;"
    done
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
    ensure_migrate
    shift
    for name in "$@"; do
      migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"
      sleep 1
    done
    ;;

  fix_dirty)
    ensure_migrate
    echo "Forcing migration state to a clean version (default 0)"
    version=${2:-0}
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$version"
    ;;

  *)
    echo "Usage: ./scripts/migrate.sh [run|run_table name|down|drop|drop_table name|clean|drop_specific name1 name2 ...|create name1 name2 ...|fix_dirty [version]]"
    ;;
esac