# pg_parse_diff

Compares PostgreSQL `nodeToString(raw_parser(...))` output to pgparser's
`nodes.NodeToString(parser.Parse(...))`.

## Prerequisites

- Build `tools/pg_parse_helper/pg_parse_helper` first.

## Usage

Single file:

```
go run tools/pg_parse_diff/main.go --file /path/to.sql
```

Directory of `.sql` files:

```
go run tools/pg_parse_diff/main.go --dir parser/pgregress/testdata/sql
```

Stdin:

```
echo "select 1;" | go run tools/pg_parse_diff/main.go
```

Custom helper path:

```
go run tools/pg_parse_diff/main.go --pg-helper /path/to/pg_parse_helper --file tools/pg_parse_diff/smoke.sql
```
