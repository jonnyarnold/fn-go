package tokeniser

import (
	"regexp"
)

// A TokenType maps a regular expression to a Tag.
type TokenType struct {
	Tag     string
	Matcher regexp.Regexp
}

func NewTokenType(tag string, matcher string) TokenType {
	compiledRegexp := regexp.MustCompile(matcher)

	return TokenType{
		Tag:     tag,
		Matcher: *compiledRegexp,
	}
}

func tt(tag string, matcher string) TokenType {
	return NewTokenType(tag, matcher)
}

var TokenTypes = []TokenType{
	tt("comment", "^#([^\n]*)"),
	tt("space", "^[\\s\n]+"),

	tt("bracket_open", "^\\("),
	tt("bracket_close", "^\\)"),

	tt("comma", "^,"),
	tt("end_statement", "^;"),

	// Reserved words/symbols
	tt("use", "^use"),
	tt("import", "^import"),
	tt("when", "^when"),

	// Infix operators
	tt("infix_operator", "^(\\+|\\-|\\*|/|\\.|=|eq|or|and)"),

	// Blocks
	tt("block_open", "^\\{"),
	tt("block_close", "^\\}"),

	// Lists
	tt("list_open", "^\\["),
	tt("list_close", "^\\]"),

	// Value literals
	tt("string", "^\"([^\"]*)\""),
	tt("number", "^([0-9]+(\\.[0-9]+)?)"),
	tt("boolean", "^(true|false)"),

	// The above regexes should catch all reserved expressions;
	// This catches the rest.
	tt("identifier", "^([^\\#\\(\\)\\,\\;\\+\\-\\*\\/\\.\\=\\|\\>\\{\\}\"\\s])+"),
}
