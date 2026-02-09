package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestCreateTempTableRelpersistence(t *testing.T) {
	input := "CREATE TEMP TABLE t (id int)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}

	if stmt.Relation.Relpersistence != 't' {
		t.Errorf("expected Relpersistence='t' for TEMP TABLE, got %q", stmt.Relation.Relpersistence)
	}
}

func TestCreateUnloggedTable(t *testing.T) {
	input := "CREATE UNLOGGED TABLE t (id int)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}

	if stmt.Relation.Relpersistence != 'u' {
		t.Errorf("expected Relpersistence='u' for UNLOGGED TABLE, got %q", stmt.Relation.Relpersistence)
	}
}

func TestCreatePermanentTable(t *testing.T) {
	input := "CREATE TABLE t (id int)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected *nodes.CreateStmt, got %T", result.Items[0])
	}

	if stmt.Relation.Relpersistence != 'p' {
		t.Errorf("expected Relpersistence='p' for permanent TABLE, got %q", stmt.Relation.Relpersistence)
	}
}

func TestCreateTempTableRelpersistenceIfNotExists(t *testing.T) {
	input := "CREATE TEMP TABLE IF NOT EXISTS t (id int)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)

	if stmt.Relation.Relpersistence != 't' {
		t.Errorf("expected Relpersistence='t' for TEMP TABLE IF NOT EXISTS, got %q", stmt.Relation.Relpersistence)
	}
	if !stmt.IfNotExists {
		t.Error("expected IfNotExists to be true")
	}
}

func TestCreateTempViewRelpersistence(t *testing.T) {
	input := "CREATE TEMP VIEW v AS SELECT 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.ViewStmt)
	if !ok {
		t.Fatalf("expected *nodes.ViewStmt, got %T", result.Items[0])
	}

	if stmt.View.Relpersistence != 't' {
		t.Errorf("expected Relpersistence='t' for TEMP VIEW, got %q", stmt.View.Relpersistence)
	}
}

func TestCreatePermanentViewRelpersistence(t *testing.T) {
	input := "CREATE VIEW v AS SELECT 1"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt, ok := result.Items[0].(*nodes.ViewStmt)
	if !ok {
		t.Fatalf("expected *nodes.ViewStmt, got %T", result.Items[0])
	}

	if stmt.View.Relpersistence != 'p' {
		t.Errorf("expected Relpersistence='p' for permanent VIEW, got %q", stmt.View.Relpersistence)
	}
}
