// Package parser provides the public API for parsing PostgreSQL SQL.
package parser

import (
	"fmt"

	"github.com/pgparser/pgparser/nodes"
)

// ParseResult contains the result of parsing.
type ParseResult struct {
	Stmts []*nodes.RawStmt
	Err   error
}

// parserLexer adapts our Lexer to the goyacc-generated pgLexer interface.
type parserLexer struct {
	lexer  *Lexer
	result *nodes.List
	err    error
}

// newParserLexer creates a new lexer adapter for the parser.
func newParserLexer(input string) *parserLexer {
	return &parserLexer{
		lexer: NewLexer(input),
	}
}

// Lex implements pgLexer.Lex
func (l *parserLexer) Lex(lval *pgSymType) int {
	tok := l.lexer.NextToken()

	// Map lexer token types to parser token types
	tokType := l.mapTokenType(tok)

	// Set semantic values based on token type
	switch tokType {
	case IDENT:
		lval.str = tok.Str
	case ICONST:
		lval.ival = tok.Ival
	case FCONST, SCONST, BCONST, XCONST:
		lval.str = tok.Str
	case Op:
		lval.str = tok.Str
	default:
		// Keywords - check if it has a string value
		if tok.Str != "" {
			lval.str = tok.Str
		}
	}

	return tokType
}

// Error implements pgLexer.Error
func (l *parserLexer) Error(s string) {
	l.err = &ParseError{
		Message:  s,
		Position: l.lexer.pos,
	}
}

// mapTokenType maps lexer token types to parser token types.
// Keywords already use parser constants (from parser.go), but
// non-keyword tokens (identifiers, literals, operators) use internal
// lex_* constants that need mapping.
func (l *parserLexer) mapTokenType(tok Token) int {
	// Handle EOF
	if tok.Type == 0 {
		return 0
	}

	// Single-character tokens map directly
	if tok.Type > 0 && tok.Type < 256 {
		return tok.Type
	}

	// Check if it's a lexer internal token (lex_* constants)
	// These start at nonKeywordTokenBase (800)
	if tok.Type >= nonKeywordTokenBase && tok.Type < nonKeywordTokenBase+100 {
		// These are non-keyword tokens from the lexer
		offset := tok.Type - nonKeywordTokenBase
		switch offset {
		case 0: // lex_ICONST
			return ICONST
		case 1: // lex_FCONST
			return FCONST
		case 2: // lex_SCONST
			return SCONST
		case 3: // lex_BCONST
			return BCONST
		case 4: // lex_XCONST
			return XCONST
		case 5: // lex_USCONST - map to SCONST for now
			return SCONST
		case 6: // lex_IDENT
			return IDENT
		case 7: // lex_UIDENT - map to IDENT for now
			return IDENT
		case 8: // lex_TYPECAST
			return TYPECAST
		case 9: // lex_DOT_DOT
			return DOT_DOT
		case 10: // lex_COLON_EQUALS
			return COLON_EQUALS
		case 11: // lex_EQUALS_GREATER
			return EQUALS_GREATER
		case 12: // lex_LESS_EQUALS
			return LESS_EQUALS
		case 13: // lex_GREATER_EQUALS
			return GREATER_EQUALS
		case 14: // lex_NOT_EQUALS
			return NOT_EQUALS
		case 15: // lex_PARAM
			return PARAM
		case 16: // lex_Op
			return Op
		}
		// Unknown lex_* token
		return 0
	}

	// All other tokens (keywords, etc.) are already parser constants
	// Just pass them through directly
	return tok.Type
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}

// Parse parses the given SQL input and returns a list of statements.
func Parse(input string) (*nodes.List, error) {
	lexer := newParserLexer(input)
	ret := pgParse(lexer)

	if lexer.err != nil {
		return nil, lexer.err
	}

	if ret != 0 {
		return nil, &ParseError{Message: fmt.Sprintf("parse error (ret=%d)", ret), Position: lexer.lexer.pos}
	}

	return lexer.result, nil
}
