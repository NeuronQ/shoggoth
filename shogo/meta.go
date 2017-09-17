package shogo

/*
ActionParser ...
*/
type ActionParser struct {
	p   Parser
	act func(interface{}) (interface{}, error)
}

func (p *ActionParser) Name() string {
	return p.p.Name()
}

func (p *ActionParser) ResultName() string {
	return p.p.ResultName()
}

func (p *ActionParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
	pos, r, err := p.p.Parse(input, pos, makeResult)
	if err != nil {
		return pos, r, err
	}
	r2, err := p.act(r)
	if err != nil {
		return pos, nil, ParsingError{
			parserName:    p.Name(),
			lastParsedPos: pos,
			errPos:        pos,
			details:       err.Error(),
		}
	}
	return pos, r2, nil
}

func Action(parser Parser, action func(interface{}) (interface{}, error)) *ActionParser {
	return &ActionParser{p: parser, act: action}
}

/*
Ref ...
*/
type RefParser struct {
	to         *Parser
	name       string
	resultName string
}

func (p *RefParser) Name() string {
	return p.name
}

func (p *RefParser) ResultName() string {
	return p.resultName
}

func (p *RefParser) Parse(input string, pos int, makeResult ASTreeMaker) (newPost int, result interface{}, rerr error) {
	return (*p.to).Parse(input, pos, makeResult)
}

func Ref(to *Parser) *RefParser {
	return &RefParser{to, "", ""}
}
