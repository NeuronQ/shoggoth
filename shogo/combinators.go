package shogo

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

/*
Seq ...
*/
type SeqParser struct {
	parsers    []Parser
	name       string
	resultName string
}

func (p *SeqParser) ResultName() string {
	return p.resultName
}

func (p *SeqParser) Name() string {
	if p.name != "" {
		return p.name
	}
	names := make([]string, len(p.parsers))
	for i, sp := range p.parsers {
		names[i] = sp.Name()
	}
	return "Seq(" + strings.Join(names, ",") + ")"
}

func (p *SeqParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	results := make([]interface{}, 0, len(p.parsers))
	lastPos := pos
	for _, sp := range p.parsers {
		nextPos, r, err := sp.Parse(input, lastPos, makeResult)
		if err != nil {
			err := err.(ParsingError)
			return pos, nil, ParsingError{
				parserName:    p.Name(),
				lastParsedPos: err.lastParsedPos,
				errPos:        err.errPos,
				expected:      err.expected,
				found:         err.found,
				details:       fmt.Sprintf("(\n    %s)", err.Error()),
			}
		}
		if r != nil {
			results = append(results, r)
		}
		lastPos = nextPos
	}
	r := interface{}(nil)
	if rt := p.ResultName(); rt != "" {
		r = makeResult(rt, results...)
	}
	return lastPos, r, nil
}

func Seq(resultName string, parsers ...Parser) *SeqParser {
	return &SeqParser{parsers, "", resultName}
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\r' || r == '\n' || r == '\t'
}

/*
SpacedSeq ...
*/
type SpacedSeqParser SeqParser

func (p *SpacedSeqParser) ResultName() string {
	return p.resultName
}

func (p *SpacedSeqParser) Name() string {
	if p.name != "" {
		return p.name
	}
	names := make([]string, len(p.parsers))
	for i, sp := range p.parsers {
		names[i] = sp.Name()
	}
	return "SpacedSeq(" + strings.Join(names, ",") + ")"
}

func (p *SpacedSeqParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	// fmt.Println("SpacedSeq.Parse")
	results := make([]interface{}, 0, len(p.parsers))
	lastPos := pos
	for _, sp := range p.parsers {
		// check if next rune is space and skip it if so
		for {
			nextRune, w := utf8.DecodeRuneInString(input[lastPos:])
			if !isSpace(nextRune) {
				break
			}
			lastPos += w
		}
		// try to parse on more, collecting results, break on failure
		nextPos, r, err := sp.Parse(input, lastPos, makeResult)
		if err != nil {
			err := err.(ParsingError)
			return pos, nil, ParsingError{
				parserName:    p.Name(),
				lastParsedPos: lastPos,
				errPos:        err.errPos,
				expected:      err.expected,
				found:         err.found,
				details:       fmt.Sprintf("(\n    %s)", err.Error()),
			}
		}
		// fmt.Printf("R::%T = %q\n", r, r)
		if r != nil {
			// fmt.Println("\tadded!")
			results = append(results, r)
		}
		lastPos = nextPos
	}
	r := interface{}(nil)
	if rt := p.ResultName(); rt != "" {
		r = makeResult(rt, results...)
	}
	return lastPos, r, nil
}

func SpacedSeq(resultName string, parsers ...Parser) *SpacedSeqParser {
	return &SpacedSeqParser{parsers, "", resultName}
}

/*
Or ...
*/
type OrParser struct {
	parsers    []Parser
	name       string
	resultName string
}

func (p *OrParser) ResultName() string {
	return p.resultName
}

func (p *OrParser) Name() string {
	if p.name != "" {
		return p.name
	}
	names := make([]string, len(p.parsers))
	for i, sp := range p.parsers {
		names[i] = sp.Name()
	}
	return "Or(" + strings.Join(names, ",") + ")"
}

func (p *OrParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	currPos := pos
	nextPos := pos
	var r interface{}
	var err error
	expected := make([]string, len(p.parsers))
	found := make([]string, len(p.parsers))
	for i, sp := range p.parsers {
		nextPos, r, err = sp.Parse(input, currPos, makeResult)

		if err == nil {
			res := interface{}(nil)
			if rt := p.ResultName(); rt != "" {
				res = makeResult(rt, r)
			}
			return nextPos, res, nil
		}

		parsingError := err.(ParsingError)
		expected[i] = parsingError.expected
		found[i] = parsingError.found

		currPos = nextPos
	}

	return pos, nil, ParsingError{
		parserName:    p.Name(),
		lastParsedPos: pos,
		errPos:        currPos,
		expected:      strings.Join(expected, "|"),
		found:         strings.Join(found, "|"),
		details:       fmt.Sprintf("(\n    %s)", err.Error()),
	}
}

func Or(resultName string, parsers ...Parser) *OrParser {
	return &OrParser{parsers, "", resultName}
}

/*
Some ...
*/
type SomeParser struct {
	p          Parser
	name       string
	resultName string
}

func (p *SomeParser) ResultName() string {
	return p.resultName
}

func (p *SomeParser) Name() string {
	if p.name != "" {
		return p.name
	}
	return "Some(" + p.p.Name() + ")"
}

func (p *SomeParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	nextPos := pos
	var r interface{}
	var err error
	results := make([]interface{}, 0, 12)
	var i int
	for ; i < 10; i++ { // TODO: replace 10 with larger number
		nextPos, r, err = p.p.Parse(input, nextPos, makeResult)
		if err != nil {
			break
		}
		if r != nil {
			results = append(results, r)
		}
	}
	if i == 0 {
		parsingError := err.(ParsingError)
		return pos, nil, ParsingError{
			parserName:    p.Name(),
			lastParsedPos: pos,
			errPos:        nextPos,
			expected:      parsingError.expected,
			found:         parsingError.found,
			details:       fmt.Sprintf("(\n    %s)", parsingError.Error()),
		}
	}
	res := interface{}(nil)
	if rt := p.ResultName(); rt != "" {
		res = makeResult(rt, results...)
	}
	return nextPos, res, nil
}

func Some(resultName string, p Parser) *SomeParser {
	return &SomeParser{p, "", resultName}
}

/*
SpacedSome ...
*/
type SpacedSomeParser SomeParser

func (p *SpacedSomeParser) ResultName() string {
	return p.resultName
}

func (p *SpacedSomeParser) Name() string {
	if p.name != "" {
		return p.name
	}
	return "SpacedSome(" + p.p.Name() + ")"
}

func (p *SpacedSomeParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	// fmt.Println("--- SpacedSome.Parse")
	nextPos := pos
	var r interface{}
	var err error
	results := make([]interface{}, 0, 12)
	var i int
	for ; i < 10; i++ { // TODO: replace 10 with larger number
		// check if next rune is space and skip it if so
		nextRune, w := utf8.DecodeRuneInString(input[nextPos:])
		if nextRune == ' ' || nextRune == '\r' || nextRune == '\n' || nextRune == '\t' {
			nextPos += w
			continue
		}
		// try to parse on more, collecting results, break on failure
		nextPos, r, err = p.p.Parse(input, nextPos, makeResult)
		if err != nil {
			break
		}
		if r != nil {
			results = append(results, r)
		}
	}
	if i == 0 {
		parsingError := err.(ParsingError)
		return pos, nil, ParsingError{
			parserName:    p.name,
			lastParsedPos: pos,
			errPos:        nextPos,
			expected:      parsingError.expected,
			found:         parsingError.found,
			details:       fmt.Sprintf("(\n    %s)", parsingError.Error()),
		}
	}
	res := interface{}(nil)
	if rt := p.ResultName(); rt != "" {
		res = makeResult(rt, results...)
	}
	return nextPos, res, nil
}

func SpacedSome(resultName string, p Parser) *SpacedSomeParser {
	return &SpacedSomeParser{p, "", resultName}
}
