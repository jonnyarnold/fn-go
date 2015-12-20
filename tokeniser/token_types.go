package tokeniser

import (
	"regexp"
)

// A TokenType maps a regular expression to a Tag.
type TokenType struct {
	Tag     string
	Matcher regexp.Regexp
	Ignore  bool
}

func NewTokenType(tag string, matcher string) TokenType {
	compiledRegexp := regexp.MustCompile(matcher)

	return TokenType{
		Tag:     tag,
		Matcher: *compiledRegexp,
		Ignore:  false,
	}
}

func IgnoredTokenType(tag string, matcher string) TokenType {
	compiledRegexp := regexp.MustCompile(matcher)

	return TokenType{
		Tag:     tag,
		Matcher: *compiledRegexp,
		Ignore:  true,
	}
}

func tt(tag string, matcher string) TokenType {
	return NewTokenType(tag, matcher)
}

func ignored(tag string, matcher string) TokenType {
	return IgnoredTokenType(tag, matcher)
}

var TokenTypes = []TokenType{
	ignored("comment", "^#([^\n]*)"),
	ignored("space", "^[\\s\n]+"),

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

	// Value literals
	tt("string", "^\"([^\"]*)\""),
	tt("number", "^([0-9]+(\\.[0-9]+)?)"),
	tt("boolean", "^(true|false)"),

	// The above regexes should catch all reserved expressions;
	// This catches the rest.
	tt("identifier", "^([^\\#\\(\\)\\,\\;\\+\\-\\*\\/\\.\\=\\|\\>\\{\\}\"\\s])+"),
}
