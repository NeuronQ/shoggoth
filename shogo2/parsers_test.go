package shogo2

import (
	"testing"

	"github.com/NeuronQ/shoggoth/sll"
	"github.com/stretchr/testify/assert"
)

func TestParsers(t *testing.T) {
	assert := assert.New(t)

	pos, r, err := Rune('x', "rune", "X")("x", 0, makeResult)
	assert.Nil(err)
	assert.Equal(1, pos)
	assert.Equal("(rune x)", r.(*sll.Element).String())

	pos, r, err = Rune('x', "letter", "Letter")("x23", 0, makeResult)
	assert.Nil(err)
	assert.Equal(1, pos)
	assert.Equal("(letter x)", r.(*sll.Element).String())

	pos, r, err = Rune('x', "", "XNothing")("x23", 0, makeResult)
	assert.Nil(err)
	assert.Equal(1, pos)
	assert.Nil(r)

	pos, r, err = Rune('x', "", "ParserX")("y", 0, makeResult)
	assert.Equal("ParsingError @ 0-1: ParserX expected \"x\", found \"y\" ", err.Error())
	assert.Equal(0, pos)
	assert.Equal(nil, r)

	pos, r, err = Rune('x', "x", "XParser")("y", 0, makeResult)
	assert.Equal("ParsingError @ 0-1: XParser expected \"x\", found \"y\" ", err.Error())
	assert.Equal(0, pos)
	assert.Nil(r)
}
