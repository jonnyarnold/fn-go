package tokeniser

import (
	"regexp"
)

// Saves some keys.
func r(regex string) regexp.Regexp {
	compiled := regexp.MustCompile(regex)
	return *compiled
}

// TokenTypes maps regular expressions to their given TokenType.
var TokenTypes = map[string]regexp.Regexp{
	"comment": r("^#([^\n]*)"),
	"space":   r("^[\\s\n]+"),

	"bracket_open":  r("^\\("),
	"bracket_close": r("^\\)"),

	"comma":         r("^,"),
	"end_statement": r("^;"),

	// Reserved words/symbols
	"use":    r("^use"),
	"import": r("^import"),
	"when":   r("^when"),

	// Infix operators
	"infix_operator": r("^(\\+|\\-|\\*|/|\\.|=|eq|or|and)"),

	// Blocks
	"block_open":  r("^\\{"),
	"block_close": r("^\\}"),

	// Lists
	"list_open":  r("^\\["),
	"list_close": r("^\\]"),

	// Value literals
	"string":  r("^\"([^\"]*)\""),
	"number":  r("^([0-9]+(?:\\.[0-9]+)?)"),
	"boolean": r("^(true|false)"),

	// The above regexes shoulds catch all reserved expressions;
	// This catches the rest.
	"identifier": r("^([^\\#\\(\\)\\,\\;\\+\\-\\*\\/\\.\\=\\|\\>\\{\\}\"\\s])+"),
}
