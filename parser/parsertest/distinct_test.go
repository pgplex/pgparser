package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestParseSelectDistinct(t *testing.T) {
	input := "SELECT DISTINCT id FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result == nil || len(result.Items) != 1 {
		t.Fatal("expected 1 statement")
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// SELECT DISTINCT should set DistinctClause to a non-nil empty list.
	// In PostgreSQL, DISTINCT (without ON) produces lcons(NIL,NIL),
	// which is an empty list indicating "distinct all".
	if stmt.DistinctClause == nil {
		t.Fatal("expected DistinctClause to be non-nil for SELECT DISTINCT")
	}
}

func TestParseSelectDistinctMultipleColumns(t *testing.T) {
	input := "SELECT DISTINCT a, b, c FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	if stmt.DistinctClause == nil {
		t.Fatal("expected DistinctClause to be non-nil for SELECT DISTINCT")
	}

	// Verify the target list still has 3 columns
	if stmt.TargetList == nil || len(stmt.TargetList.Items) != 3 {
		t.Fatalf("expected 3 target columns, got %v", stmt.TargetList)
	}
}

func TestParseSelectDistinctOn(t *testing.T) {
	input := "SELECT DISTINCT ON (a) a, b FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// DISTINCT ON should set DistinctClause to a list containing the ON expressions.
	if stmt.DistinctClause == nil {
		t.Fatal("expected DistinctClause to be non-nil for SELECT DISTINCT ON")
	}
	if len(stmt.DistinctClause.Items) != 1 {
		t.Fatalf("expected 1 DISTINCT ON expression, got %d", len(stmt.DistinctClause.Items))
	}

	// The expression should be a ColumnRef for "a"
	colRef, ok := stmt.DistinctClause.Items[0].(*nodes.ColumnRef)
	if !ok {
		t.Fatalf("expected *nodes.ColumnRef in DistinctClause, got %T", stmt.DistinctClause.Items[0])
	}
	colStr := colRef.Fields.Items[0].(*nodes.String)
	if colStr.Str != "a" {
		t.Errorf("expected DISTINCT ON column 'a', got %q", colStr.Str)
	}
}

func TestParseSelectWithoutDistinct(t *testing.T) {
	input := "SELECT id FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *nodes.SelectStmt, got %T", result.Items[0])
	}

	// A plain SELECT (no DISTINCT) should have nil DistinctClause.
	if stmt.DistinctClause != nil {
		t.Errorf("expected DistinctClause to be nil for plain SELECT, got %v", stmt.DistinctClause)
	}
}
