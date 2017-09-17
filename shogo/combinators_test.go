package shogo

import (
	"testing"

	"github.com/NeuronQ/shoggoth/sll"
	"github.com/stretchr/testify/assert"
)

func TestCombinators_Seq(t *testing.T) {
	pos, r, err := Seq("number", Rune('0'), Rune('1')).Parse("01", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 2, pos)
	assert.Equal(t, "(number (rune 0) (rune 1))", r.(*sll.Element).String())

	pos, r, err = Seq("number", Rune('a'), Rune('b')).Parse("ba", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)
}

func TestCombinators_SpacedSeq(t *testing.T) {
	pos, r, err := SpacedSeq("number", Rune('0'), Rune('1')).Parse("0  1", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 4, pos)
	assert.Equal(t, "(number (rune 0) (rune 1))", r.(*sll.Element).String())

	pos, r, err = SpacedSeq("number", Rune('0'), Rune('1')).Parse("  0  1", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 6, pos)
	assert.Equal(t, "(number (rune 0) (rune 1))", r.(*sll.Element).String())

	pos, r, err = SpacedSeq("number", Rune('a'), Rune('b')).Parse("ba", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)
}

func TestCombinators_Or(t *testing.T) {
	pos, r, err := Or("number", Rune('0'), Rune('1')).Parse("01", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 1, pos)
	assert.Equal(t, "(number (rune 0))", r.(*sll.Element).String())

	pos, r, err = Or("number", Rune('0'), Rune('1')).Parse("10", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 1, pos)
	assert.Equal(t, "(number (rune 1))", r.(*sll.Element).String())

	pos, r, err = Or("number", Rune('a'), Rune('b'), Rune('c')).Parse("c", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 1, pos)
	assert.Equal(t, "(number (rune c))", r.(*sll.Element).String())

	pos, r, err = Or("number", Rune('a'), Rune('b'), Rune('c')).Parse("3", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)
}

func TestCombinators_Some(t *testing.T) {
	pos, r, err := Some("xs", Rune('X')).Parse("XXX", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 3, pos)
	assert.Equal(t, "(xs (rune X) (rune X) (rune X))", r.(*sll.Element).String())

	pos, r, err = Some("xs", Rune('X')).Parse("Y32", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)

	pos, r, err = Some("xs", Rune('X')).Parse(" XXX", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)
}

func TestCombinators_SpacedSome(t *testing.T) {
	pos, r, err := Some("xs", Rune('X')).Parse("XXX", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 3, pos)
	assert.Equal(t, "(xs (rune X) (rune X) (rune X))", r.(*sll.Element).String())

	pos, r, err = SpacedSome("xs", Rune('X')).Parse("  X\tXX ", 0, makeResult)
	assert.Nil(t, err)
	assert.Equal(t, 7, pos)
	assert.Equal(t, "(xs (rune X) (rune X) (rune X))", r.(*sll.Element).String())

	pos, r, err = SpacedSome("xs", Rune('X')).Parse("Y32", 0, makeResult)
	assert.NotNil(t, err)
	assert.Equal(t, 0, pos)
	assert.Equal(t, nil, r)
}
