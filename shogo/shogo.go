package shogo

import "fmt"

// Parser is the interface for all parsers
type Parser interface {
	Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error)
	Name() string
	ResultName() string
}

// ASTreeMaker is the signature of a function taking a parser name and some parsing result values (of any kind) and returning whatever AST-style structure is needed.
type ASTreeMaker func(parserName string, values ...interface{}) interface{}

// ParsingError is the error structure we use for capturing and displaying detailed error messages
type ParsingError struct {
	parserName    string
	lastParsedPos int
	errPos        int
	expected      string
	found         string
	details       string
}

func (e ParsingError) Error() string {
	return fmt.Sprintf("ParsingError @ %d-%d: %s expected \"%s\", found \"%s\" %s",
		e.lastParsedPos, e.errPos, e.parserName, e.expected, e.found, e.details)
}
