package shogo

import (
	"testing"

	"github.com/NeuronQ/shoggoth/sll"
	"github.com/stretchr/testify/assert"
)

func TestParsers(t *testing.T) {
	assert := assert.New(t)

	pos, r, err := Rune('x').Parse("x", 0, makeResult)
	assert.Equal(nil, err)
	assert.Equal(1, pos)
	assert.Equal("(rune x)", r.(*sll.Element).String())

	pos, r, err = Rune('x').Parse("x23", 0, makeResult)
	assert.Equal(nil, err)
	assert.Equal(1, pos)
	assert.Equal("(rune x)", r.(*sll.Element).String())

	pos, r, err = Rune('x').Parse("y", 0, makeResult)
	assert.Equal("ParsingError @ 0-1: Rune(x) expected \"x\", found \"y\" ", err.Error())
	assert.Equal(0, pos)
	assert.Equal(nil, r)
}
