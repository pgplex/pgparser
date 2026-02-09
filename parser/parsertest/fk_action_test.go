package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

// findFKConstraint finds a FOREIGN KEY constraint in table elements.
func findFKConstraint(t *testing.T, stmt *nodes.CreateStmt) *nodes.Constraint {
	t.Helper()
	for _, elt := range stmt.TableElts.Items {
		switch v := elt.(type) {
		case *nodes.Constraint:
			if v.Contype == nodes.CONSTR_FOREIGN {
				return v
			}
		case *nodes.ColumnDef:
			if v.Constraints != nil {
				for _, c := range v.Constraints.Items {
					if cc, ok := c.(*nodes.Constraint); ok && cc.Contype == nodes.CONSTR_FOREIGN {
						return cc
					}
				}
			}
		}
	}
	return nil
}

func TestForeignKeyOnDeleteSetNullColumns(t *testing.T) {
	input := "CREATE TABLE t (a int, b int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL (a, b))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	fk := findFKConstraint(t, stmt)
	if fk == nil {
		t.Fatal("expected a FOREIGN KEY constraint")
	}

	if fk.FkDelaction != 'n' {
		t.Errorf("expected FkDelaction='n' (SET NULL), got %q", fk.FkDelaction)
	}

	if fk.FkDelsetcols == nil {
		t.Fatal("expected FkDelsetcols to be set for ON DELETE SET NULL (col_list)")
	}
	if len(fk.FkDelsetcols.Items) != 2 {
		t.Fatalf("expected 2 columns in FkDelsetcols, got %d", len(fk.FkDelsetcols.Items))
	}
}

func TestForeignKeyOnDeleteSetDefaultColumns(t *testing.T) {
	input := "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET DEFAULT (a))"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	fk := findFKConstraint(t, stmt)
	if fk == nil {
		t.Fatal("expected a FOREIGN KEY constraint")
	}

	if fk.FkDelaction != 'd' {
		t.Errorf("expected FkDelaction='d' (SET DEFAULT), got %q", fk.FkDelaction)
	}

	if fk.FkDelsetcols == nil {
		t.Fatal("expected FkDelsetcols to be set for ON DELETE SET DEFAULT (col_list)")
	}
	if len(fk.FkDelsetcols.Items) != 1 {
		t.Fatalf("expected 1 column in FkDelsetcols, got %d", len(fk.FkDelsetcols.Items))
	}
}

func TestForeignKeyOnDeleteSetNullNoColumns(t *testing.T) {
	input := "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES p(id) ON DELETE SET NULL)"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.CreateStmt)
	fk := findFKConstraint(t, stmt)
	if fk == nil {
		t.Fatal("expected a FOREIGN KEY constraint")
	}

	if fk.FkDelaction != 'n' {
		t.Errorf("expected FkDelaction='n' (SET NULL), got %q", fk.FkDelaction)
	}

	if fk.FkDelsetcols != nil {
		t.Errorf("expected FkDelsetcols to be nil without column list, got %v", fk.FkDelsetcols)
	}
}
