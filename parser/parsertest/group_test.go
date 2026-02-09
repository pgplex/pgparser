package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestParseGroupByDistinct(t *testing.T) {
	input := "SELECT a, count(*) FROM t GROUP BY DISTINCT a"

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

	// GROUP BY DISTINCT should set GroupDistinct to true
	if !stmt.GroupDistinct {
		t.Fatal("expected GroupDistinct to be true for GROUP BY DISTINCT")
	}

	// GroupClause should contain the column "a"
	if stmt.GroupClause == nil || len(stmt.GroupClause.Items) != 1 {
		t.Fatalf("expected 1 GROUP BY item, got %v", stmt.GroupClause)
	}
}

func TestParseGroupByWithoutDistinct(t *testing.T) {
	input := "SELECT a, count(*) FROM t GROUP BY a"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// Plain GROUP BY should have GroupDistinct = false
	if stmt.GroupDistinct {
		t.Error("expected GroupDistinct to be false for plain GROUP BY")
	}

	if stmt.GroupClause == nil || len(stmt.GroupClause.Items) != 1 {
		t.Fatalf("expected 1 GROUP BY item, got %v", stmt.GroupClause)
	}
}

func TestParseGroupByAll(t *testing.T) {
	input := "SELECT a, count(*) FROM t GROUP BY ALL a"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)

	// GROUP BY ALL should have GroupDistinct = false
	if stmt.GroupDistinct {
		t.Error("expected GroupDistinct to be false for GROUP BY ALL")
	}

	if stmt.GroupClause == nil || len(stmt.GroupClause.Items) != 1 {
		t.Fatalf("expected 1 GROUP BY item, got %v", stmt.GroupClause)
	}
}
