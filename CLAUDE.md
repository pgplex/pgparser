# pgparser - Pure Go PostgreSQL Parser

## Project Overview

pgparser is a pure Go implementation of PostgreSQL's SQL parser, providing 100% compatible parse tree representation. It translates PostgreSQL's gram.y (Bison grammar) into Go using goyacc.

**Target PostgreSQL Version**: 17.7 (REL_17_STABLE)

## Key Directories

- `parser/` - SQL parser (goyacc-based, translated from PostgreSQL gram.y)
- `nodes/` - Node type definitions (parse tree types, enums, serialization)
- `cmd/pgsema-gen/` - Code generator for keywords and node types

## PostgreSQL Source References

All translations are based on `~/Github/postgres` (REL_17_STABLE branch):

- Grammar: `src/backend/parser/gram.y`
- Lexer: `src/backend/parser/scan.l`
- Keywords: `src/include/parser/kwlist.h`
- Parse nodes: `src/include/nodes/parsenodes.h`
- Primitive nodes: `src/include/nodes/primnodes.h`

## Development Workflow

### Build

```bash
go build ./...
```

### Run Tests

```bash
go test ./...
```

### Regenerate Parser

```bash
make generate-parser
```

### Regenerate Node Types

```bash
make generate-nodes
```

## Translation Guidelines

### goyacc Translation (from Bison)

1. `%locations` - Not supported; add location fields to %union manually
2. `%parse-param` / `%lex-param` - Embed state in Lexer struct
3. `makeNode(X)` -> `&nodes.X{}`
4. `list_make1(x)` -> `[]nodes.Node{x}`
5. `lappend(list, x)` -> `append(list, x)`
6. `$1`, `$2` -> Need explicit type assertions: `$1.(*nodes.SelectStmt)`

### Node Serialization

Output format must match PostgreSQL's `nodeToString()` exactly:
```
{QUERY :commandType 1 :querySource 0 :canSetTag true ...}
```

Field order must match PostgreSQL's output order.
