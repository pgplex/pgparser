# gram.y Conflict Resolution Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reduce gram.y parser conflicts from 1,481 (437 S/R + 1,044 R/R) to ~65 by fixing structural translation errors from PostgreSQL's Bison grammar.

**Architecture:** The conflicts stem from 3 root causes: (1) missing lexer lookahead tokens (`FORMAT_LA`, `WITHOUT_LA`, incomplete `WITH_LA`), (2) duplicate/redundant grammar productions, and (3) inherent LALR(1) ambiguities. We fix categories 1 and 2 which account for ~1,416 conflicts. Category 3 (~65 conflicts) are benign — resolved correctly by goyacc defaults.

**Tech Stack:** goyacc, Go, PostgreSQL gram.y reference (`~/Github/postgres/src/backend/parser/gram.y`)

**Reference files:**
- PostgreSQL grammar: `~/Github/postgres/src/backend/parser/gram.y`
- PostgreSQL lexer lookahead: `~/Github/postgres/src/backend/parser/parser.c:195-251`
- pgparser grammar: `parser/gram.y`
- pgparser lexer: `parser/parse.go`
- pgparser tests: `parser/parser_test.go`

---

### Task 1: Add `FORMAT_LA` and `WITHOUT_LA` token declarations to gram.y

**Files:**
- Modify: `parser/gram.y:39-41` (token declarations)

**Step 1: Add token declarations**

At `parser/gram.y:41`, after `NULLS_LA`, add the two new lookahead tokens:

```yacc
%token         NOT_LA    /* lookahead token for NOT LIKE etc */
%token         WITH_LA   /* lookahead token for WITH TIME ZONE */
%token         NULLS_LA  /* lookahead token for NULLS FIRST/LAST */
%token         FORMAT_LA /* lookahead token for FORMAT JSON */
%token         WITHOUT_LA /* lookahead token for WITHOUT TIME ZONE */
```

**Step 2: Regenerate parser and verify tokens compile**

Run: `cd parser && go generate`
Expected: Compiles without errors (conflicts will still exist but new tokens are available)

**Step 3: Commit**

```bash
git add parser/gram.y
git commit -m "feat(parser): declare FORMAT_LA and WITHOUT_LA lookahead tokens"
```

---

### Task 2: Add `FORMAT_LA` and `WITHOUT_LA` to lexer lookahead logic

**Files:**
- Modify: `parser/parse.go:107-129` (applyLookahead function)

**Step 1: Write the failing test**

Add to `parser/parser_test.go`:

```go
func TestLookaheadFormatLA(t *testing.T) {
	// FORMAT followed by JSON should become FORMAT_LA
	input := "COPY t TO STDOUT (FORMAT JSON)"
	_, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error for FORMAT JSON: %v", err)
	}
}

func TestLookaheadWithoutLA(t *testing.T) {
	// WITHOUT followed by TIME should become WITHOUT_LA
	input := "SELECT CAST('2024-01-01' AS TIMESTAMP WITHOUT TIME ZONE)"
	_, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error for WITHOUT TIME ZONE: %v", err)
	}
}

func TestLookaheadWithOrdinality(t *testing.T) {
	// WITH followed by ORDINALITY should become WITH_LA
	input := "SELECT * FROM generate_series(1,3) WITH ORDINALITY"
	_, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error for WITH ORDINALITY: %v", err)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./parser/ -run "TestLookahead(Format|Without|WithOrdinality)" -v`
Expected: FAIL (at least some of these should fail due to parse errors from token ambiguity)

**Step 3: Implement lookahead changes**

In `parser/parse.go`, modify the `applyLookahead` function to add 3 new cases:

```go
func (l *parserLexer) applyLookahead(curToken, nextToken int) int {
	switch curToken {
	case NOT:
		switch nextToken {
		case BETWEEN, IN_P, LIKE, ILIKE, SIMILAR:
			return NOT_LA
		}
	case WITH:
		// Replace WITH by WITH_LA if followed by TIME or ORDINALITY
		switch nextToken {
		case TIME, ORDINALITY:
			return WITH_LA
		}
	case NULLS_P:
		switch nextToken {
		case FIRST_P, LAST_P, DISTINCT, NOT:
			return NULLS_LA
		}
	case FORMAT:
		// Replace FORMAT by FORMAT_LA if followed by JSON
		switch nextToken {
		case JSON:
			return FORMAT_LA
		}
	case WITHOUT:
		// Replace WITHOUT by WITHOUT_LA if followed by TIME
		switch nextToken {
		case TIME:
			return WITHOUT_LA
		}
	}
	return curToken
}
```

**Step 4: Run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./parser/ -run "TestLookahead" -v`
Expected: Tests may still fail because grammar rules haven't been updated yet. That's OK — we'll fix in the following tasks.

**Step 5: Commit**

```bash
git add parser/parse.go parser/parser_test.go
git commit -m "feat(parser): add FORMAT_LA, WITHOUT_LA, WITH+ORDINALITY lookahead"
```

---

### Task 3: Update `with_clause` and `opt_with` to use `WITH_LA` (~393 S/R conflicts)

This is the single biggest conflict source (391 S/R in state 884 + 2 S/R in opt_ordinality).

**Files:**
- Modify: `parser/gram.y:6223-6226` (opt_with)
- Modify: `parser/gram.y:7876-7891` (with_clause)
- Modify: `parser/gram.y:8380-8383` (opt_ordinality)

**Step 1: Update `with_clause` to accept `WITH_LA`**

At `parser/gram.y:7876`, add `WITH_LA cte_list` alternative (matching PostgreSQL `gram.y:12874`):

```yacc
with_clause:
	WITH cte_list
		{
			$$ = &nodes.WithClause{
				Ctes:      $2,
				Recursive: false,
			}
		}
	| WITH_LA cte_list
		{
			$$ = &nodes.WithClause{
				Ctes:      $2,
				Recursive: false,
			}
		}
	| WITH RECURSIVE cte_list
		{
			$$ = &nodes.WithClause{
				Ctes:      $3,
				Recursive: true,
			}
		}
	;
```

**Step 2: Update `opt_with` to accept `WITH_LA`**

At `parser/gram.y:6223`:

```yacc
opt_with:
	WITH      {}
	| WITH_LA {}
	| /* EMPTY */ {}
	;
```

**Step 3: Update `opt_ordinality` to use `WITH_LA`**

At `parser/gram.y:8380`:

```yacc
opt_ordinality:
	WITH_LA ORDINALITY		{ $$ = true }
	| /* EMPTY */		{ $$ = false }
	;
```

**Step 4: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: All existing tests pass. Check y.output for reduced conflict count.

**Step 5: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: S/R count drops by ~393 (from 437 to ~44)

**Step 6: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output
git commit -m "fix(parser): add WITH_LA to with_clause/opt_with/opt_ordinality (-393 conflicts)"
```

---

### Task 4: Update `opt_timezone` to use `WITHOUT_LA` (~6 S/R conflicts)

**Files:**
- Modify: `parser/gram.y:11671-11675` (opt_timezone)

**Step 1: Update opt_timezone**

At `parser/gram.y:11671`:

```yacc
opt_timezone:
	WITH_LA TIME ZONE		{ $$ = true }
	| WITHOUT_LA TIME ZONE		{ $$ = false }
	| /* EMPTY */			{ $$ = false }
	;
```

Change `WITHOUT TIME ZONE` to `WITHOUT_LA TIME ZONE`.

**Step 2: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: Tests pass. S/R count drops by ~6.

**Step 3: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: 6 fewer S/R conflicts

**Step 4: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output
git commit -m "fix(parser): use WITHOUT_LA in opt_timezone (-6 conflicts)"
```

---

### Task 5: Rewrite `utility_option_name` with `FORMAT_LA` (~863 R/R conflicts)

This is the second biggest conflict source. PostgreSQL uses `NonReservedWord` + `FORMAT_LA`, while pgparser lists individual keywords that overlap with `ColId`.

**Files:**
- Modify: `parser/gram.y:5201-5212` (utility_option_name)
- Modify: `parser/gram.y:10569-10570` (json_format_clause — FORMAT → FORMAT_LA)

**Step 1: Rewrite `utility_option_name`**

At `parser/gram.y:5201`, replace the current production with PostgreSQL's design. Use `NonReservedWord` (which pgparser already has at line 5280) instead of `ColId` plus redundant keyword alternatives:

```yacc
utility_option_name:
	NonReservedWord  { $$ = $1 }
	| ANALYZE        { $$ = "analyze" }
	| FORMAT_LA      { $$ = "format" }
	;
```

Rationale:
- `NonReservedWord` = `IDENT | unreserved_keyword | col_name_keyword | type_func_name_keyword` — covers FORMAT, FORCE, VERBOSE, FULL, FREEZE, CONCURRENTLY (all unreserved or type_func_name keywords)
- `ANALYZE` is a reserved keyword, not reachable via `NonReservedWord`, so it needs explicit listing
- `FORMAT_LA` is the lookahead-replaced token when FORMAT precedes JSON, needed since `FORMAT_LA` is not a keyword and not in `NonReservedWord`
- Remove NULL_P and DEFAULT — these are reserved keywords. If needed, add them explicitly. Check PostgreSQL: it does NOT list them. They are handled through `opt_boolean_or_string` or other paths. **Actually check**: look at PostgreSQL's `utility_option_name` — it only has `NonReservedWord | analyze_keyword | FORMAT_LA`. So NULL_P and DEFAULT should NOT be here.

**Step 2: Update `json_format_clause` to use `FORMAT_LA`**

Search for all occurrences of `FORMAT JSON` in gram.y and replace with `FORMAT_LA JSON`. The json_format_clause at line 10569:

```yacc
json_format_clause:
	FORMAT_LA JSON
		{
			$$ = &nodes.JsonFormat{
				FormatType: nodes.JS_FORMAT_JSON,
				Encoding:   nodes.JS_ENC_DEFAULT,
			}
		}
	| FORMAT_LA JSON ENCODING name
		{
			...
		}
	;
```

Also check `json_output_type_clause` and any other rule that uses `FORMAT JSON` — all must use `FORMAT_LA JSON`.

**Step 3: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: Tests pass. R/R count drops dramatically (~852 from FORMAT/FORCE + ~11 S/R from FORMAT alias contexts).

**Step 4: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: ~863 fewer conflicts

**Step 5: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output
git commit -m "fix(parser): rewrite utility_option_name with NonReservedWord+FORMAT_LA (-863 conflicts)"
```

---

### Task 6: Remove duplicate CAST production from `c_expr` (~112 R/R conflicts)

**Files:**
- Modify: `parser/gram.y:9503-9510` (c_expr CAST rule)

**Step 1: Remove the duplicate**

Delete lines 9503-9510 from `c_expr`. The CAST production already exists in `func_expr_common_subexpr` and is reachable via `c_expr -> func_expr -> func_expr_common_subexpr`.

Before (c_expr):
```yacc
	| ARRAY array_expr
		{
			$$ = $2
		}
	| CAST '(' a_expr AS Typename ')'        // <-- DELETE THIS BLOCK
		{                                      // <-- DELETE
			$$ = &nodes.TypeCast{              // <-- DELETE
				Arg:      $3,                  // <-- DELETE
				TypeName: $5,                  // <-- DELETE
				Location: -1,                  // <-- DELETE
			}                                  // <-- DELETE
		}                                      // <-- DELETE
	| explicit_row
```

After:
```yacc
	| ARRAY array_expr
		{
			$$ = $2
		}
	| explicit_row
```

**Step 2: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: Tests pass. Verify CAST still works via the existing test suite.

**Step 3: Write a specific regression test for CAST**

Add to `parser/parser_test.go`:

```go
func TestParseCast(t *testing.T) {
	tests := []string{
		"SELECT CAST(1 AS INTEGER)",
		"SELECT CAST('hello' AS TEXT)",
		"SELECT CAST(x AS NUMERIC(10,2)) FROM t",
	}
	for _, sql := range tests {
		_, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse error for %q: %v", sql, err)
		}
	}
}
```

**Step 4: Run the new test**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./parser/ -run TestParseCast -v`
Expected: PASS

**Step 5: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: ~112 fewer R/R conflicts

**Step 6: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output parser/parser_test.go
git commit -m "fix(parser): remove duplicate CAST in c_expr, keep in func_expr_common_subexpr (-112 conflicts)"
```

---

### Task 7: Remove redundant `select_no_parens` rules (~18 R/R conflicts)

**Files:**
- Modify: `parser/gram.y:7681-7686,7712-7716` (redundant select rules)

**Step 1: Remove the two redundant rules**

Delete `select_clause opt_sort_clause select_limit` (line 7681-7686) — it is a subset of `select_clause opt_sort_clause select_limit opt_for_locking_clause` (line 7693-7698) when `opt_for_locking_clause` is empty.

Also delete `with_clause select_clause opt_sort_clause select_limit` (line 7712-7716) — same reason.

After removal, `select_no_parens` should have exactly 8 alternatives (matching PostgreSQL's 8):

```yacc
select_no_parens:
	simple_select
		{ $$ = $1 }
	| select_clause sort_clause
		{
			n := $1.(*nodes.SelectStmt)
			n.SortClause = $2
			$$ = n
		}
	| select_clause opt_sort_clause for_locking_clause opt_select_limit
		{
			n := $1.(*nodes.SelectStmt)
			insertSelectOptions(n, $2, $3, $4, nil)
			$$ = n
		}
	| select_clause opt_sort_clause select_limit opt_for_locking_clause
		{
			n := $1.(*nodes.SelectStmt)
			insertSelectOptions(n, $2, $4, $3, nil)
			$$ = n
		}
	| with_clause select_clause
		{
			n := $2.(*nodes.SelectStmt)
			n.WithClause = $1.(*nodes.WithClause)
			$$ = n
		}
	| with_clause select_clause sort_clause
		{
			n := $2.(*nodes.SelectStmt)
			n.WithClause = $1.(*nodes.WithClause)
			n.SortClause = $3
			$$ = n
		}
	| with_clause select_clause opt_sort_clause for_locking_clause opt_select_limit
		{
			n := $2.(*nodes.SelectStmt)
			insertSelectOptions(n, $3, $4, $5, $1.(*nodes.WithClause))
			$$ = n
		}
	| with_clause select_clause opt_sort_clause select_limit opt_for_locking_clause
		{
			n := $2.(*nodes.SelectStmt)
			insertSelectOptions(n, $3, $5, $4, $1.(*nodes.WithClause))
			$$ = n
		}
	;
```

**Step 2: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: PASS

**Step 3: Add regression test for SELECT with LIMIT**

```go
func TestParseSelectLimit(t *testing.T) {
	tests := []string{
		"SELECT * FROM t LIMIT 10",
		"SELECT * FROM t ORDER BY a LIMIT 10",
		"SELECT * FROM t ORDER BY a LIMIT 10 OFFSET 5",
		"SELECT * FROM t LIMIT 10 FOR UPDATE",
		"SELECT * FROM t FOR UPDATE LIMIT 10",
		"SELECT * FROM t ORDER BY a LIMIT 10 FOR UPDATE",
	}
	for _, sql := range tests {
		_, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse error for %q: %v", sql, err)
		}
	}
}
```

**Step 4: Run test and verify**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./parser/ -run TestParseSelectLimit -v`
Expected: PASS

**Step 5: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: ~18 fewer R/R conflicts

**Step 6: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output parser/parser_test.go
git commit -m "fix(parser): remove redundant select_no_parens rules (-18 conflicts)"
```

---

### Task 8: Remove redundant keyword alternatives from `privilege` (~24 R/R conflicts)

**Files:**
- Modify: `parser/gram.y:6038-6059` (privilege production)

**Step 1: Remove redundant unreserved keyword alternatives**

INSERT, UPDATE, DELETE_P, TRIGGER, EXECUTE, TRUNCATE are all `unreserved_keyword` and already reachable via `ColId`. Remove them and add `ALTER SYSTEM_P` (present in PostgreSQL but missing in pgparser).

```yacc
privilege:
	SELECT opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "select", Cols: $2} }
	| REFERENCES opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "references", Cols: $2} }
	| CREATE opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: "create", Cols: $2} }
	| ALTER SYSTEM_P
		{ $$ = &nodes.AccessPriv{PrivName: "alter system"} }
	| ColId opt_column_list
		{ $$ = &nodes.AccessPriv{PrivName: $1, Cols: $2} }
	;
```

Only SELECT, REFERENCES, CREATE are reserved keywords not reachable via ColId, so they need explicit alternatives. All other privilege names (insert, update, delete, trigger, execute, truncate, usage, etc.) are unreserved and handled by the `ColId opt_column_list` catch-all.

**Step 2: Regenerate parser and run tests**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && cd parser && go generate && cd .. && go test ./parser/ -v`
Expected: PASS

**Step 3: Add regression test for GRANT**

```go
func TestParseGrant(t *testing.T) {
	tests := []string{
		"GRANT SELECT ON TABLE t TO role1",
		"GRANT INSERT ON TABLE t TO role1",
		"GRANT UPDATE ON TABLE t TO role1",
		"GRANT DELETE ON TABLE t TO role1",
		"GRANT ALL ON TABLE t TO role1",
		"GRANT SELECT, INSERT, UPDATE ON TABLE t TO role1",
	}
	for _, sql := range tests {
		_, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse error for %q: %v", sql, err)
		}
	}
}
```

**Step 4: Run test**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./parser/ -run TestParseGrant -v`
Expected: PASS

**Step 5: Verify conflict reduction**

Run: `tail -1 parser/y.output`
Expected: ~24 fewer R/R conflicts

**Step 6: Commit**

```bash
git add parser/gram.y parser/parser.go parser/y.output parser/parser_test.go
git commit -m "fix(parser): remove redundant privilege keyword alternatives (-24 conflicts)"
```

---

### Task 9: Final verification and conflict count baseline

**Step 1: Run full test suite**

Run: `cd /private/var/folders/hz/x716z9xj68g7njffd675mljr0000gn/T/vibe-kanban/worktrees/17c0-gram-y/pgparser && go test ./... -v`
Expected: All tests pass

**Step 2: Check final conflict count**

Run: `tail -1 parser/y.output`
Expected: approximately `~33 shift/reduce, ~32 reduce/reduce conflicts reported` (total ~65)

**Step 3: Document remaining conflicts**

The remaining ~65 conflicts are benign LALR(1) ambiguities that goyacc resolves correctly via default behavior (shift-prefer for S/R, first-rule-prefer for R/R). These match patterns in PostgreSQL's own grammar that Bison also resolves implicitly. They include:

- `[` in indirection/array_bounds contexts (8 S/R) — shift is correct
- ColId vs ColLabel in function_with_argtypes (18 R/R) — first rule is correct
- DROP SUBSCRIPTION name vs name_list (8 mixed) — shift is correct
- any_name vs function_with_argtypes (10 R/R) — first rule is correct
- ARRAY `[` dimension (2 S/R) — shift is correct
- VALUES `(` (1 S/R) — shift is correct
- GENERATED AS IDENTITY (2 S/R) — shift is correct
- ConstraintAttributeSpec (3 S/R) — shift is correct
- OptTempTableName + FORMAT (3 R/R) — resolved by FORMAT_LA now, may go to 0
- substr_list (2 R/R) — first rule is correct
- BEGIN ATOMIC END (1 R/R) — first rule is correct
- NESTED PATH (1 S/R) — shift is correct

**Step 4: Commit**

```bash
git add -A
git commit -m "docs: document remaining benign parser conflicts (~65)"
```

---

## Summary of Expected Impact

| Task | Target | Conflicts Eliminated | Cumulative |
|------|--------|---------------------|------------|
| 1-2 | Add FORMAT_LA, WITHOUT_LA tokens + lexer | 0 (preparation) | 1,481 |
| 3 | WITH_LA in with_clause/opt_with/opt_ordinality | ~393 | ~1,088 |
| 4 | WITHOUT_LA in opt_timezone | ~6 | ~1,082 |
| 5 | FORMAT_LA in utility_option_name + json_format_clause | ~863 | ~219 |
| 6 | Remove duplicate CAST in c_expr | ~112 | ~107 |
| 7 | Remove redundant select_no_parens rules | ~18 | ~89 |
| 8 | Remove redundant privilege keywords | ~24 | ~65 |
| **Total eliminated** | | **~1,416** | **~65 remaining** |
