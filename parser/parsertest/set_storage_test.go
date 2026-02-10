package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestAlterColumnSetStoragePlain(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN c SET STORAGE PLAIN"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	alt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := alt.Cmds.Items[0].(*nodes.AlterTableCmd)

	if cmd.Def == nil {
		t.Fatal("expected Def to be set with storage mode String node")
	}
	str, ok := cmd.Def.(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String for Def, got %T", cmd.Def)
	}
	if str.Str != "plain" {
		t.Errorf("expected storage mode 'plain', got %q", str.Str)
	}
}

func TestAlterColumnSetStorageExternal(t *testing.T) {
	input := "ALTER TABLE t ALTER c SET STORAGE EXTERNAL"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	alt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := alt.Cmds.Items[0].(*nodes.AlterTableCmd)

	if cmd.Def == nil {
		t.Fatal("expected Def to be set")
	}
	str := cmd.Def.(*nodes.String)
	if str.Str != "external" {
		t.Errorf("expected storage mode 'external', got %q", str.Str)
	}
}

func TestAlterColumnSetStorageDefault(t *testing.T) {
	input := "ALTER TABLE t ALTER COLUMN c SET STORAGE DEFAULT"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	alt := result.Items[0].(*nodes.AlterTableStmt)
	cmd := alt.Cmds.Items[0].(*nodes.AlterTableCmd)

	if cmd.Def == nil {
		t.Fatal("expected Def to be set")
	}
	str := cmd.Def.(*nodes.String)
	if str.Str != "default" {
		t.Errorf("expected storage mode 'default', got %q", str.Str)
	}
}
