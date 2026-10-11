#!/usr/bin/env bash
set -euo pipefail

# Base directory for migrations
MIGRATE_DIR="./sql/postgres/migrations"

# Migration name passed as argument, default to "new_migration"
NAME=${1:-new_migration}

# Ensure directory exists
mkdir -p "$MIGRATE_DIR"

# Find the highest existing migration number
LAST=$(ls -1 "$MIGRATE_DIR" 2>/dev/null | sed -E 's/([0-9]+).*/\1/' | sort -n | tail -1 2>/dev/null)
if [ -z "$LAST" ]; then
  LAST=0
fi
NEXT=$((LAST+1))
NUM=$(printf "%06d" $NEXT)

# Create the up and down files
touch "$MIGRATE_DIR/$NUM.$NAME.up.sql"
touch "$MIGRATE_DIR/$NUM.$NAME.down.sql"

# Write template content
cat > "$MIGRATE_DIR/$NUM.$NAME.up.sql" <<'EOF'
-- Create migration: replace with actual SQL
EOF

cat > "$MIGRATE_DIR/$NUM.$NAME.down.sql" <<'EOF'
-- Revert migration: replace with DROP statements
EOF

printf "Created migration %s.$NAME.up.sql and %s.$NAME.down.sql\n" "$NUM" "$NUM"
