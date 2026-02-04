package grammar

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytebase/pgparser/nodes"
	"github.com/bytebase/pgparser/parser"
	"github.com/bytebase/pgparser/parser/pgregress"
)

type expectations map[string][]string

var fileExpectations = expectations{
	"phase1_nodes.sql": {
		"SQLValueFunction",
		"CoalesceExpr",
		"MinMaxExpr",
	},
	"phase5_plan.sql": {
		"JsonTableColumn",
		"JsonFormat",
		"CreateStatsStmt",
		"StatsElem",
		"DefineStmt",
		"CreateFunctionStmt",
		"TypeName",
	},
	"phase6_extract.sql": {
		"CreateFunctionStmt",
		"RuleStmt",
	},
}

func TestGrammarFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("sql", "*.sql"))
	if err != nil {
		t.Fatalf("glob fixtures: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no grammar fixtures found")
	}

	for _, path := range files {
		path := path
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			stmts := pgregress.ExtractStatements(name, content)
			if len(stmts) == 0 {
				t.Fatalf("no statements in %s", name)
			}

			found := make(map[string]struct{})
			for _, stmt := range stmts {
				if stmt.SQL == "" {
					continue
				}
				expectErr, sql := extractExpectedError(stmt.SQL)
				if strings.TrimSpace(sql) == "" {
					continue
				}
				result, err := parser.Parse(sql)
				if expectErr != "" {
					if err == nil {
						t.Fatalf("expected error %q in %s line %d", expectErr, name, stmt.StartLine)
					}
					if !strings.Contains(err.Error(), expectErr) {
						t.Fatalf("expected error %q in %s line %d, got %q", expectErr, name, stmt.StartLine, err.Error())
					}
					continue
				}
				if err != nil {
					t.Fatalf("parse error in %s line %d: %v", name, stmt.StartLine, err)
				}
				if result == nil {
					t.Fatalf("parse returned nil for %s line %d", name, stmt.StartLine)
				}
				collectNodeTypes(result, found)
			}

			expect, ok := fileExpectations[name]
			if !ok {
				return
			}

			for _, typeName := range expect {
				if _, ok := found[typeName]; !ok {
					t.Fatalf("expected node type %s in %s", typeName, name)
				}
			}
		})
	}
}

var nodeType = reflect.TypeOf((*nodes.Node)(nil)).Elem()

func extractExpectedError(sql string) (string, string) {
	lines := strings.Split(sql, "\n")
	cleaned := make([]string, 0, len(lines))
	expectErr := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- expect-error:") {
			if expectErr == "" {
				expectErr = strings.TrimSpace(strings.TrimPrefix(trimmed, "-- expect-error:"))
			}
			continue
		}
		cleaned = append(cleaned, line)
	}
	return expectErr, strings.Join(cleaned, "\n")
}

func collectNodeTypes(n nodes.Node, out map[string]struct{}) {
	seen := make(map[uintptr]struct{})
	walk(reflect.ValueOf(n), out, seen)
}

func walk(v reflect.Value, out map[string]struct{}, seen map[uintptr]struct{}) {
	if !v.IsValid() {
		return
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		ptr := v.Pointer()
		if ptr != 0 {
			if _, ok := seen[ptr]; ok {
				return
			}
			seen[ptr] = struct{}{}
		}
		if v.Type().Implements(nodeType) {
			out[v.Elem().Type().Name()] = struct{}{}
		}
		walk(v.Elem(), out, seen)
		return
	}

	if v.CanAddr() && v.Addr().Type().Implements(nodeType) {
		out[v.Type().Name()] = struct{}{}
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			walk(v.Field(i), out, seen)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			walk(v.Index(i), out, seen)
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			walk(iter.Value(), out, seen)
		}
	}
}
