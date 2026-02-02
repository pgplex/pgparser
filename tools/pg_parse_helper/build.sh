#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PG_SRC="${PG_SRC:-$HOME/Github/postgres}"
PG_CONFIG="${PG_CONFIG:-$PG_SRC/src/bin/pg_config/pg_config}"

if [[ ! -d "$PG_SRC" ]]; then
  echo "PG_SRC not found: $PG_SRC" >&2
  exit 1
fi

if [[ ! -f "$PG_SRC/src/include/pg_config.h" ]]; then
  echo "Missing pg_config.h. Configure and build Postgres source first:" >&2
  echo "  cd \"$PG_SRC\" && ./configure && make -j && make -C src/backend libpostgres.a" >&2
  exit 1
fi

if [[ ! -f "$PG_SRC/src/backend/libpostgres.a" ]]; then
  echo "Missing libpostgres.a. Build it first:" >&2
  echo "  make -C \"$PG_SRC/src/backend\" libpostgres.a" >&2
  exit 1
fi

PLATFORM="$(uname -s | tr '[:upper:]' '[:lower:]')"
PORT_DIR="$PG_SRC/src/include/port/$PLATFORM"

INCLUDES=(
  "-I$PG_SRC/src/include"
  "-I$PG_SRC/src/include/port"
  "-I$PG_SRC/src/backend"
)

if [[ -d "$PORT_DIR" ]]; then
  INCLUDES+=("-I$PORT_DIR")
fi

LDFLAGS=()
LIBS=()
if [[ -x "$PG_CONFIG" ]]; then
  # Use pg_config to pick up system libs if available.
  LDFLAGS+=($("$PG_CONFIG" --ldflags))
  LIBS+=($("$PG_CONFIG" --libs))
fi

cc -O2 "${INCLUDES[@]}" \
  "$ROOT_DIR/pg_parse_helper.c" \
  "$PG_SRC/src/backend/libpostgres.a" \
  "$PG_SRC/src/common/libpgcommon.a" \
  "$PG_SRC/src/port/libpgport.a" \
  "${LDFLAGS[@]}" "${LIBS[@]}" \
  -o "$ROOT_DIR/pg_parse_helper"

echo "Built $ROOT_DIR/pg_parse_helper"
