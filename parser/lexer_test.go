package parser

import (
	"testing"
)

func TestLexerKeywords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		str      string
	}{
		{"SELECT", SELECT, "select"},
		{"select", SELECT, "select"},
		{"FROM", FROM, "from"},
		{"WHERE", WHERE, "where"},
		{"AND", AND, "and"},
		{"OR", OR, "or"},
		{"CREATE", CREATE, "create"},
		{"TABLE", TABLE, "table"},
		{"INSERT", INSERT, "insert"},
		{"UPDATE", UPDATE, "update"},
		{"DELETE", DELETE_P, "delete"},
		{"NULL", NULL_P, "null"},
		{"TRUE", TRUE_P, "true"},
		{"FALSE", FALSE_P, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.expected {
				t.Errorf("expected token type %d, got %d", tt.expected, tok.Type)
			}
			if tok.Str != tt.str {
				t.Errorf("expected str %q, got %q", tt.str, tok.Str)
			}
		})
	}
}

func TestLexerIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		str      string
	}{
		{"foo", lex_IDENT, "foo"},
		{"Foo", lex_IDENT, "foo"},
		{"FOO", lex_IDENT, "foo"},
		{"foo_bar", lex_IDENT, "foo_bar"},
		{"_foo", lex_IDENT, "_foo"},
		{"foo123", lex_IDENT, "foo123"},
		{"foo$bar", lex_IDENT, "foo$bar"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.expected {
				t.Errorf("expected token type %d (IDENT), got %d", tt.expected, tok.Type)
			}
			if tok.Str != tt.str {
				t.Errorf("expected str %q, got %q", tt.str, tok.Str)
			}
		})
	}
}

func TestLexerDelimitedIdentifiers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foo"`, "foo"},
		{`"Foo Bar"`, "Foo Bar"},
		{`"SELECT"`, "SELECT"},
		{`"foo""bar"`, `foo"bar`}, // escaped quote
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != lex_IDENT {
				t.Errorf("expected token type IDENT, got %d", tok.Type)
			}
			if tok.Str != tt.expected {
				t.Errorf("expected str %q, got %q", tt.expected, tok.Str)
			}
		})
	}
}

func TestLexerNumbers(t *testing.T) {
	tests := []struct {
		input    string
		tokType  int
		ival     int64
		str      string
	}{
		{"0", lex_ICONST, 0, "0"},
		{"42", lex_ICONST, 42, "42"},
		{"123456", lex_ICONST, 123456, "123456"},
		{"1_000_000", lex_ICONST, 1000000, "1000000"},
		{"0x10", lex_ICONST, 16, "0x10"},
		{"0xFF", lex_ICONST, 255, "0xFF"},
		{"0o17", lex_ICONST, 15, "0o17"},
		{"0b1010", lex_ICONST, 10, "0b1010"},
		{"3.14", lex_FCONST, 0, "3.14"},
		{".5", lex_FCONST, 0, ".5"},
		{"1e10", lex_FCONST, 0, "1e10"},
		{"1.5e-3", lex_FCONST, 0, "1.5e-3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.tokType {
				t.Errorf("expected token type %d, got %d", tt.tokType, tok.Type)
			}
			if tt.tokType == lex_ICONST && tok.Ival != tt.ival {
				t.Errorf("expected ival %d, got %d", tt.ival, tok.Ival)
			}
			if tt.tokType == lex_FCONST && tok.Str != tt.str {
				t.Errorf("expected str %q, got %q", tt.str, tok.Str)
			}
		})
	}
}

func TestLexerStrings(t *testing.T) {
	tests := []struct {
		input    string
		tokType  int
		expected string
	}{
		{"'hello'", lex_SCONST, "hello"},
		{"'hello world'", lex_SCONST, "hello world"},
		{"'it''s'", lex_SCONST, "it's"},      // escaped quote
		{"''", lex_SCONST, ""},               // empty string
		{"E'hello\\nworld'", lex_SCONST, "hello\nworld"}, // escape string
		{"E'tab\\there'", lex_SCONST, "tab\there"},
		{"B'101'", lex_BCONST, "b101"},       // bit string
		{"X'FF'", lex_XCONST, "xFF"},         // hex string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.tokType {
				t.Errorf("expected token type %d, got %d", tt.tokType, tok.Type)
			}
			if tok.Str != tt.expected {
				t.Errorf("expected str %q, got %q", tt.expected, tok.Str)
			}
		})
	}
}

func TestLexerDollarQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"$$hello$$", "hello"},
		{"$$hello world$$", "hello world"},
		{"$tag$hello$tag$", "hello"},
		{"$foo$some 'text' here$foo$", "some 'text' here"},
		{"$$line1\nline2$$", "line1\nline2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != lex_SCONST {
				t.Errorf("expected token type lex_SCONST, got %d", tok.Type)
			}
			if tok.Str != tt.expected {
				t.Errorf("expected str %q, got %q", tt.expected, tok.Str)
			}
		})
	}
}

func TestLexerOperators(t *testing.T) {
	tests := []struct {
		input   string
		tokType int
		str     string
	}{
		{"::", lex_TYPECAST, "::"},
		{"..", lex_DOT_DOT, ".."},
		{":=", lex_COLON_EQUALS, ":="},
		{"=>", lex_EQUALS_GREATER, "=>"},
		{"<=", lex_LESS_EQUALS, "<="},
		{">=", lex_GREATER_EQUALS, ">="},
		{"<>", lex_NOT_EQUALS, "<>"},
		{"!=", lex_NOT_EQUALS, "!="},
		{"+", int('+'), "+"},
		{"-", int('-'), "-"},
		{"*", int('*'), "*"},
		{"/", int('/'), "/"},
		{"@", lex_Op, "@"},
		{"||", lex_Op, "||"},
		{"->", lex_Op, "->"},
		{"->>", lex_Op, "->>"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.tokType {
				t.Errorf("expected token type %d, got %d", tt.tokType, tok.Type)
			}
			if tok.Str != tt.str {
				t.Errorf("expected str %q, got %q", tt.str, tok.Str)
			}
		})
	}
}

func TestLexerParameters(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"$1", 1},
		{"$2", 2},
		{"$10", 10},
		{"$123", 123},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != lex_PARAM {
				t.Errorf("expected token type lex_PARAM, got %d", tok.Type)
			}
			if tok.Ival != tt.expected {
				t.Errorf("expected ival %d, got %d", tt.expected, tok.Ival)
			}
		})
	}
}

func TestLexerComments(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		str      string
	}{
		{"-- comment\nSELECT", SELECT, "select"},
		{"/* comment */ SELECT", SELECT, "select"},
		{"/* multi\nline */ SELECT", SELECT, "select"},
		{"/* nested /* comment */ */ SELECT", SELECT, "select"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.expected {
				t.Errorf("expected token type %d, got %d", tt.expected, tok.Type)
			}
			if tok.Str != tt.str {
				t.Errorf("expected str %q, got %q", tt.str, tok.Str)
			}
		})
	}
}

func TestLexerWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"  SELECT", SELECT},
		{"\tSELECT", SELECT},
		{"\nSELECT", SELECT},
		{"\r\nSELECT", SELECT},
		{"  \t\n  SELECT", SELECT},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tok := lexer.NextToken()
			if tok.Type != tt.expected {
				t.Errorf("expected token type %d, got %d", tt.expected, tok.Type)
			}
		})
	}
}

func TestLexerSQLStatements(t *testing.T) {
	input := "SELECT * FROM users WHERE id = 1"
	expected := []struct {
		tokType int
		str     string
	}{
		{SELECT, "select"},
		{int('*'), "*"},
		{FROM, "from"},
		{lex_IDENT, "users"},
		{WHERE, "where"},
		{lex_IDENT, "id"},
		{int('='), "="},
		{lex_ICONST, "1"},
		{lex_EOF, ""},
	}

	lexer := NewLexer(input)
	for i, tt := range expected {
		tok := lexer.NextToken()
		if tok.Type != tt.tokType {
			t.Errorf("token %d: expected type %d, got %d", i, tt.tokType, tok.Type)
		}
		if tt.tokType != lex_EOF && tt.tokType != lex_ICONST && tok.Str != tt.str {
			t.Errorf("token %d: expected str %q, got %q", i, tt.str, tok.Str)
		}
	}
}

func TestLexerCreateTable(t *testing.T) {
	input := `CREATE TABLE users (
		id INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email TEXT
	)`

	lexer := NewLexer(input)
	tokens := []Token{}
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lex_EOF {
			break
		}
	}

	// Just verify we got through without errors
	if lexer.Err != nil {
		t.Errorf("unexpected error: %v", lexer.Err)
	}

	// Verify some key tokens
	if tokens[0].Type != CREATE {
		t.Errorf("expected CREATE, got %d", tokens[0].Type)
	}
	if tokens[1].Type != TABLE {
		t.Errorf("expected TABLE, got %d", tokens[1].Type)
	}
}

func TestLexerLocation(t *testing.T) {
	input := "SELECT foo"
	lexer := NewLexer(input)

	tok := lexer.NextToken()
	if tok.Loc != 0 {
		t.Errorf("expected SELECT location 0, got %d", tok.Loc)
	}

	tok = lexer.NextToken()
	if tok.Loc != 7 {
		t.Errorf("expected foo location 7, got %d", tok.Loc)
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		{"'unterminated", "unterminated quoted string"},
		{`"unterminated`, "unterminated quoted identifier"},
		{"/* unterminated", "unterminated /* comment"},
		{"$$unterminated", "unterminated dollar-quoted string"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			for {
				tok := lexer.NextToken()
				if tok.Type == lex_EOF {
					break
				}
			}
			if lexer.Err == nil {
				t.Errorf("expected error containing %q, got nil", tt.err)
			} else if lexer.Err.Error() != tt.err {
				t.Errorf("expected error %q, got %q", tt.err, lexer.Err.Error())
			}
		})
	}
}
