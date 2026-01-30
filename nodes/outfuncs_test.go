package nodes

import (
	"strings"
	"testing"
)

func TestNodeToString_SelectStmt(t *testing.T) {
	// Create a simple SELECT statement: SELECT 1
	stmt := &SelectStmt{
		TargetList: &List{
			Items: []Node{
				&ResTarget{
					Val: &A_Const{
						Val: &Integer{Ival: 1},
					},
				},
			},
		},
	}

	result := NodeToString(stmt)

	if !strings.Contains(result, "SELECTSTMT") {
		t.Errorf("expected output to contain SELECTSTMT, got: %s", result)
	}
	if !strings.Contains(result, ":targetList") {
		t.Errorf("expected output to contain :targetList, got: %s", result)
	}
}

func TestNodeToString_SelectWithFrom(t *testing.T) {
	// Create: SELECT * FROM users
	stmt := &SelectStmt{
		TargetList: &List{
			Items: []Node{
				&ResTarget{
					Val: &ColumnRef{
						Fields: &List{
							Items: []Node{&A_Star{}},
						},
					},
				},
			},
		},
		FromClause: &List{
			Items: []Node{
				&RangeVar{
					Relname: "users",
					Inh:     true,
				},
			},
		},
	}

	result := NodeToString(stmt)

	if !strings.Contains(result, "SELECTSTMT") {
		t.Errorf("expected output to contain SELECTSTMT, got: %s", result)
	}
	if !strings.Contains(result, ":fromClause") {
		t.Errorf("expected output to contain :fromClause, got: %s", result)
	}
	if !strings.Contains(result, "RANGEVAR") {
		t.Errorf("expected output to contain RANGEVAR, got: %s", result)
	}
	if !strings.Contains(result, `"users"`) {
		t.Errorf("expected output to contain \"users\", got: %s", result)
	}
}

func TestNodeToString_Integer(t *testing.T) {
	n := &Integer{Ival: 42}
	result := NodeToString(n)
	if result != "42" {
		t.Errorf("expected 42, got: %s", result)
	}
}

func TestNodeToString_String(t *testing.T) {
	n := &String{Str: "hello"}
	result := NodeToString(n)
	if result != `"hello"` {
		t.Errorf(`expected "hello", got: %s`, result)
	}
}

func TestNodeToString_List(t *testing.T) {
	n := &List{
		Items: []Node{
			&Integer{Ival: 1},
			&Integer{Ival: 2},
			&Integer{Ival: 3},
		},
	}
	result := NodeToString(n)
	if result != "(1 2 3)" {
		t.Errorf("expected (1 2 3), got: %s", result)
	}
}

func TestNodeToString_Nil(t *testing.T) {
	result := NodeToString(nil)
	if result != "<>" {
		t.Errorf("expected <>, got: %s", result)
	}
}

func TestNodeToString_CreateStmt(t *testing.T) {
	// CREATE TABLE users (id INTEGER)
	stmt := &CreateStmt{
		Relation: &RangeVar{
			Relname: "users",
			Inh:     true,
		},
		TableElts: &List{
			Items: []Node{
				&ColumnDef{
					Colname: "id",
					TypeName: &TypeName{
						Names: &List{
							Items: []Node{
								&String{Str: "pg_catalog"},
								&String{Str: "int4"},
							},
						},
					},
				},
			},
		},
	}

	result := NodeToString(stmt)

	if !strings.Contains(result, "CREATESTMT") {
		t.Errorf("expected output to contain CREATESTMT, got: %s", result)
	}
	if !strings.Contains(result, ":tableElts") {
		t.Errorf("expected output to contain :tableElts, got: %s", result)
	}
	if !strings.Contains(result, "COLUMNDEF") {
		t.Errorf("expected output to contain COLUMNDEF, got: %s", result)
	}
}

func TestNodeToString_BoolExpr(t *testing.T) {
	// a AND b
	expr := &BoolExpr{
		Boolop: AND_EXPR,
		Args: &List{
			Items: []Node{
				&ColumnRef{
					Fields: &List{
						Items: []Node{&String{Str: "a"}},
					},
				},
				&ColumnRef{
					Fields: &List{
						Items: []Node{&String{Str: "b"}},
					},
				},
			},
		},
	}

	result := NodeToString(expr)

	if !strings.Contains(result, "BOOLEXPR") {
		t.Errorf("expected output to contain BOOLEXPR, got: %s", result)
	}
	if !strings.Contains(result, ":boolop 0") { // AND_EXPR = 0
		t.Errorf("expected output to contain :boolop 0, got: %s", result)
	}
}

func TestNodeToString_EscapeString(t *testing.T) {
	n := &String{Str: "hello\nworld\"test\\end"}
	result := NodeToString(n)
	expected := `"hello\nworld\"test\\end"`
	if result != expected {
		t.Errorf("expected %s, got: %s", expected, result)
	}
}
