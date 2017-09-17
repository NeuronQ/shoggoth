package shogo2

import (
	"testing"

	"github.com/NeuronQ/shoggoth/sll"
	"github.com/stretchr/testify/assert"
)

func dummyParser(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
	if input[pos:] == "dummy" {
		return pos + 5, makeResult("dummyR", "dummyV"), nil
	} else {
		return pos, nil, ParsingError{
			parserName:    "DummyParser",
			lastParsedPos: pos,
			errPos:        pos,
		}
	}
}

func TestMeta(t *testing.T) {
	assert := assert.New(t)

	pos, r, err := Action(dummyParser, func(r interface{}) (interface{}, error) {
		return r.(*sll.Element).Get(0)
	})("dummy", 0, makeResult)
	assert.Nil(err)
	assert.Equal(5, pos)
	assert.Equal("dummyR", r.(string))

	pos, r, err = Action(dummyParser, func(r interface{}) (interface{}, error) {
		return r.(*sll.Element).Get(0)
	})("no dummy", 0, makeResult)
	assert.Equal("ParsingError @ 0-0: DummyParser expected \"\", found \"\" ", err.Error())
	assert.Equal(0, pos)
	assert.Nil(r)

	var varDummyParser Parser
	varDummyParser = dummyParser
	pos, r, err = Ref(&varDummyParser)("dummy", 0, makeResult)
	assert.Nil(err)
	assert.Equal(5, pos)
	assert.Equal("(dummyR dummyV)", r.(*sll.Element).String())

	pos, r, err = Ref(&varDummyParser)("no dummy", 0, makeResult)
	assert.Equal("ParsingError @ 0-0: DummyParser expected \"\", found \"\" ", err.Error())
	assert.Equal(0, pos)
	assert.Nil(r)
}
