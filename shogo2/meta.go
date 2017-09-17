package shogo2

import "fmt"

func Action(p Parser, act func(interface{}) (interface{}, error)) Parser {
	return func(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
		pos, r, err := p(input, pos, makeResult)
		if err != nil {
			return pos, r, err
		}
		r2, err := act(r)
		if err != nil {
			return pos, nil, ParsingError{
				parserName:    fmt.Sprintf("Action(%q, %q)", p, act),
				lastParsedPos: pos,
				errPos:        pos,
				details:       err.Error(),
			}
		}
		return pos, r2, nil
	}
}

func Ref(p *Parser) Parser {
	return func(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
		return (*p)(input, pos, makeResult)
	}
}
