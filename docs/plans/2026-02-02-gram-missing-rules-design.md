# Grammar Completion (Missing Rules) Design

Date: 2026-02-02

## Goal
Close all missing grammar rule alternatives to reach PostgreSQL syntax coverage, while keeping AST construction and nodeToString parity on track. Work is ordered by dependency to minimize rework.

## Phases and Dependencies

### M1: Infrastructure (Nodes / %union / Helpers)
- Add missing node types required by upcoming grammar branches.
- Extend `%union` with new node pointer types.
- Add Go helper constructors used by grammar actions.
- Establish a new test harness under `tests/grammar/` (no changes to `parser/pgregress`).

### M2: Core Expressions
- `func_expr_common_subexpr` (SQL value functions, CAST, EXTRACT, SUBSTRING, TRIM, TREAT, etc.).
- `c_expr` extensions (subquery postfix, indirection, ROW/VALUES variants).
- `a_expr` / `b_expr` missing branches (OPERATOR, ANY/ALL, IS [NOT] DISTINCT, BETWEEN, COLLATE).
- row/array/indirection support.

### M3: Query Skeleton
- `table_ref` / `joined_table` / `from_clause` / lateral / tablesample / rows from.
- `with_clause` / CTE list.
- window clause/definitions/frame clause.
- group clause + grouping sets + rollup/cube.
- order/limit/locking clause.

### M4: DDL + Advanced Features
- Core DDL (CREATE/ALTER/DROP) and common admin syntax (COPY/GRANT/REVOKE/COMMENT).
- JSON/XML/partitioning/CTE extensions (search/cycle).

## Testing Strategy
- New tests live in `tests/grammar/` with SQL fixtures under `tests/grammar/sql/`.
- Each new node type must be triggered by at least one SQL statement in tests.
- Tests focus on parse success and basic node-type assertions (not nodeToString parity).

## Milestones and Acceptance

### M1 Acceptance
- `go build ./...` and `go test ./...` pass (ignoring pre-existing failures).
- New tests pass.
- Each new node type has at least one coverage statement.

### M2/M3/M4 Acceptance
- Phase-specific SQL fixtures pass.
- Added grammar branches each have at least one test case.

