package shogo

import (
	"unicode/utf8"
)

type RuneParser struct {
	r          rune
	name       string
	resultName string
}

func (p *RuneParser) Name() string {
	return p.name
}

func (p *RuneParser) ResultName() string {
	return p.resultName
}

func (p *RuneParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
	nextRune, w := utf8.DecodeRuneInString(input[pos:])
	if nextRune != p.r {
		return pos, nil, ParsingError{
			parserName:    p.name,
			lastParsedPos: pos,
			errPos:        pos + w,
			expected:      string(p.r),
			found:         input[pos : pos+w],
		}
	}
	r := interface{}(nil)
	if p.resultName != "" {
		r = makeResult(p.resultName, string(p.r))
	}
	return pos + w, r, nil
}

func Rune(r rune) *RuneParser {
	return &RuneParser{
		r:          r,
		name:       "Rune(" + string(r) + ")",
		resultName: "rune",
	}
}
