package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestTableCheckNotValid(t *testing.T) {
	input := "CREATE TABLE t (x int, CHECK (x > 0) NOT VALID)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	// Find the CHECK constraint in table elements
	var chk *nodes.Constraint
	for _, elt := range stmt.TableElts.Items {
		if c, ok := elt.(*nodes.Constraint); ok && c.Contype == nodes.CONSTR_CHECK {
			chk = c
			break
		}
	}
	if chk == nil {
		t.Fatal("expected a CHECK constraint")
	}

	if !chk.SkipValidation {
		t.Error("expected SkipValidation to be true for NOT VALID")
	}
	if chk.InitiallyValid {
		t.Error("expected InitiallyValid to be false for NOT VALID")
	}
}

func TestTableCheckValid(t *testing.T) {
	input := "CREATE TABLE t (x int, CHECK (x > 0))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	var chk *nodes.Constraint
	for _, elt := range stmt.TableElts.Items {
		if c, ok := elt.(*nodes.Constraint); ok && c.Contype == nodes.CONSTR_CHECK {
			chk = c
			break
		}
	}
	if chk == nil {
		t.Fatal("expected a CHECK constraint")
	}

	if chk.SkipValidation {
		t.Error("expected SkipValidation to be false without NOT VALID")
	}
	if !chk.InitiallyValid {
		t.Error("expected InitiallyValid to be true without NOT VALID")
	}
}

func TestColumnCheckInitiallyValid(t *testing.T) {
	// Column-level CHECK constraint (does not support NOT VALID, but should
	// still have InitiallyValid=true to match PostgreSQL defaults).
	input := "CREATE TABLE t (x int CHECK (x > 0))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	// Column-level constraints are inside ColumnDef.Constraints
	var chk *nodes.Constraint
	for _, elt := range stmt.TableElts.Items {
		if cd, ok := elt.(*nodes.ColumnDef); ok && cd.Constraints != nil {
			for _, c := range cd.Constraints.Items {
				if cc, ok := c.(*nodes.Constraint); ok && cc.Contype == nodes.CONSTR_CHECK {
					chk = cc
					break
				}
			}
		}
	}
	if chk == nil {
		t.Fatal("expected a column-level CHECK constraint")
	}

	if !chk.InitiallyValid {
		t.Error("expected InitiallyValid to be true for column-level CHECK (matching PostgreSQL)")
	}
}
