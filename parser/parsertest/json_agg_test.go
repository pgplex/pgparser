package parsertest

import (
	"testing"

	"github.com/pgplex/pgparser/nodes"
	"github.com/pgplex/pgparser/parser"
)

func TestJsonObjectAggFilter(t *testing.T) {
	input := "SELECT JSON_OBJECTAGG(k VALUE v) FILTER (WHERE k IS NOT NULL) FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	agg, ok := rt.Val.(*nodes.JsonObjectAgg)
	if !ok {
		t.Fatalf("expected *nodes.JsonObjectAgg, got %T", rt.Val)
	}

	if agg.Constructor == nil {
		t.Fatal("expected Constructor to be set")
	}
	if agg.Constructor.Agg_filter == nil {
		t.Fatal("expected Agg_filter to be set for FILTER clause")
	}
}

func TestJsonArrayAggOver(t *testing.T) {
	input := "SELECT JSON_ARRAYAGG(v) OVER (PARTITION BY grp) FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	agg, ok := rt.Val.(*nodes.JsonArrayAgg)
	if !ok {
		t.Fatalf("expected *nodes.JsonArrayAgg, got %T", rt.Val)
	}

	if agg.Constructor == nil {
		t.Fatal("expected Constructor to be set")
	}
	if agg.Constructor.Over == nil {
		t.Fatal("expected Over to be set for OVER clause")
	}
}

func TestJsonObjectAggNoFilterNoOver(t *testing.T) {
	input := "SELECT JSON_OBJECTAGG(k VALUE v) FROM t"

	result, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	stmt := result.Items[0].(*nodes.SelectStmt)
	rt := stmt.TargetList.Items[0].(*nodes.ResTarget)
	agg, ok := rt.Val.(*nodes.JsonObjectAgg)
	if !ok {
		t.Fatalf("expected *nodes.JsonObjectAgg, got %T", rt.Val)
	}

	if agg.Constructor.Agg_filter != nil {
		t.Error("expected Agg_filter to be nil without FILTER clause")
	}
	if agg.Constructor.Over != nil {
		t.Error("expected Over to be nil without OVER clause")
	}
}
