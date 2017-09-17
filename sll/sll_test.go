package sll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var E = New

func TestSll_Len(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(0, E().Len())
	assert.Equal(1, E(nil).Len())
	assert.Equal(1, E(42).Len())
	assert.Equal(2, E(11, 22).Len())
	assert.Equal(2, E("foo", 10).Len())
	l := E(
		11,
		22,
		E(
			E("go", []int{42, 042, 0x42}),
			-3,
		),
		33,
		"ok",
	)
	assert.Equal(5, l.Len())
}

func TestSll_String(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("()", E().String())
	assert.Equal("(<nil>)", E(nil).String())
	assert.Equal("(42)", E(42).String())
	assert.Equal("(11 22)", E(11, 22).String())
	assert.Equal("(foo 10)", E("foo", 10).String())
	l := E(
		11,
		22,
		E(
			E("go", []int{42, 042, 0x42}),
			-3,
		),
		33,
		"ok",
	)
	assert.Equal("(11 22 ((go [42 34 66]) -3) 33 ok)", l.String())
}

func TestSll_Dump(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("[]", E().Dump())
	assert.Equal("[<nil>::<nil>]", E(nil).Dump())
	assert.Equal("[42::int]", E(42).Dump())
	assert.Equal("[11::int 22::int]", E(11, 22).Dump())
	assert.Equal("[foo::string 10::int]", E("foo", 10).Dump())
	l := E(
		11,
		22,
		E(
			E("go", []int{42, 042, 0x42}),
			-3,
		),
		33,
		"ok",
	)
	assert.Equal("[11::int 22::int [[go::string [42 34 66]::[]int] -3::int] 33::int ok::string]", l.Dump())
}

func TestSll_Get(t *testing.T) {
	assert := assert.New(t)

	r, err := E(nil).Get(0)
	assert.Nil(err)
	assert.Equal(nil, r)

	r, err = E(nil).Get(1)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = E(42).Get(0)
	assert.Nil(err)
	assert.Equal(42, r)

	r, err = E(42).Get(1)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = E(42).Get(12)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = E(11, "howdy").Get(1)
	assert.Nil(err)
	assert.Equal("howdy", r)

	r, err = E(11, "howdy").Get(2)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = E(11, "howdy", -3.14).Get(2)
	assert.Nil(err)
	assert.Equal(-3.14, r)

	r, err = E(11, "howdy", -3.14).Get(3)
	assert.NotNil(err)
	assert.Nil(r)

	l := E(
		11,
		22,
		E(
			E("go", []int{42, 042, 0x42}),
			-3,
		),
		33,
		"ok",
	)
	r, err = l.Get(0)
	assert.Nil(err)
	assert.Equal(11, r)

	r, err = l.Get(4)
	assert.Nil(err)
	assert.Equal("ok", r)

	r, err = l.Get(5)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = l.Get(2, 1)
	assert.Nil(err)
	assert.Equal(-3, r)

	r, err = l.Get(2, 2)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = l.Get(2, 0, 0)
	assert.Nil(err)
	assert.Equal("go", r)

	r, err = l.Get(2, 0, 2)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = l.Get(2, 0, 1)
	assert.Nil(err)
	assert.Equal(42, r.([]int)[0])

	r, err = l.Get(2, 0, 2, 0)
	assert.NotNil(err)
	assert.Nil(r)

	r, err = l.Get(2, 0, 2, 13, 0)
	assert.NotNil(err)
	assert.Nil(r)
}

func TestSll_GetList(t *testing.T) {
	assert := assert.New(t)

	// (20 30 40 50) -> (30)
	l := E(20, 30, 40, 50)
	r, err := l.GetList([]int{1})
	assert.Nil(err)
	assert.Equal("(30)", r.String())

	// (20 30 40 50) -> (30 40 50)
	l = E(20, 30, 40, 50)
	r, err = l.GetList([]int{1}, []int{2}, []int{3})
	assert.Nil(err)
	assert.Equal("(30 40 50)", r.String())

	// (20 30 40 50) -> (30 40 50)
	l = E(20, 30, 40, 50)
	r, err = l.GetList([]int{1, 4})
	assert.NotNil(err)
	assert.Nil(r)

	// (op (rune +)) -> (op +)
	l = E("op", E("rune", "x"))
	r, err = l.GetList([]int{0}, []int{1, 1})
	assert.Nil(err)
	assert.Equal("(op x)", r.String())
}
