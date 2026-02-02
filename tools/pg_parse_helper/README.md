# pg_parse_helper

Minimal C helper that calls PostgreSQL's `raw_parser()` and prints
`nodeToString()` for the raw parse tree list. Intended as a baseline for
pgparser diff tooling.

## Prerequisites

- A configured and built PostgreSQL source tree (default: `~/Github/postgres`)
- `libpostgres.a` present in `src/backend/`

If the source tree is not configured yet:

```
cd ~/Github/postgres
./configure
make -j
make -C src/backend libpostgres.a
```

## Build

```
PG_SRC=~/Github/postgres tools/pg_parse_helper/build.sh
```

This produces `tools/pg_parse_helper/pg_parse_helper`.

## Usage

```
tools/pg_parse_helper/pg_parse_helper --file /path/to/sql.sql
```

Or from stdin:

```
echo "select 1;" | tools/pg_parse_helper/pg_parse_helper
```

Modes:

- `--mode=default` (default)
- `--mode=type_name`
- `--mode=plpgsql_expr`
- `--mode=plpgsql_assign1`
- `--mode=plpgsql_assign2`
- `--mode=plpgsql_assign3`
