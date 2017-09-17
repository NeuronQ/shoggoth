/*
Package sll implement the simples imaginalble single-linked list, with useful methods for accessing deelply nested list (aka "trees"), and decent string representations.

It's mostly inspired by Lisp's Cons data structure used to represent trees (like "abstract suntax trees" aka ASTs) like nested lists of lists, also called S-expressions.

Due to both its conceptual origins and the desire for the simplest possible implementation, an "element" and a "list" are the same thing: an element is the "head" of the list
starting from it. Nevertheless, this is meant to be used only as a list, so the "cons pair" behavior from the Lisp world is not supported as this can only cause confusion in a language like Go.

It's `Get` is loosely inspired by Clojure's `get-in` but adapted to be more Go-idiomatic.
*/
package sll

import (
	"errors"
	"fmt"
	"strings"
)

// Element represents the (head) element of a singly linked list
type Element struct {
	Value interface{}
	Next  *Element
}

// New creates a new singly linked list with the given args as values
func New(vals ...interface{}) *Element {
	var r *Element
	for i := len(vals) - 1; i >= 0; i-- {
		r = &Element{Value: vals[i], Next: r}
		// fmt.Printf("r = %+v\n", r)
	}
	return r
}

// Len gives the lengths of a (possibly nested) lists
func (e *Element) Len() int {
	if e == nil {
		return 0
	}
	sz := 1
	for p := e; p.Next != nil; p = p.Next {
		// fmt.Println(p.Value)
		sz++
		if sz > 10 {
			break
		}
	}
	return sz
}

// String gives an S-expression style (parens and values) string representation
func (e *Element) String() string {
	words := make([]string, 0, 3*e.Len())
	for p := e; p != nil; p = p.Next {
		words = append(words, fmt.Sprint(p.Value))
	}
	return "(" + strings.Join(words, " ") + ")"
}

// Dump returns a more debugging-friendly string representation
func (e *Element) Dump() string {
	words := make([]string, 0, 3*e.Len())
	for p := e; p != nil; p = p.Next {
		switch v := p.Value.(type) {
		case *Element:
			words = append(words, v.Dump())
		default:
			words = append(words, fmt.Sprintf("%v::%T", v, v))
		}
	}
	return "[" + strings.Join(words, " ") + "]"
}

// Get fetches an element e[ks[0]]...[ks[n]] from a (possible nested) list
func (e *Element) Get(ks ...int) (interface{}, error) {
	if e == nil {
		return nil, errors.New("empty list")
	}
	v := e
	iMax := len(ks) - 1
	for i, k := range ks {
		for j := 0; j < k; j++ {
			if v.Next == nil {
				return nil, fmt.Errorf("invalid index: %d (arg %d), max here is %d", k, i, j)
			}
			v = v.Next
		}
		if i < iMax {
			subTree, ok := v.Value.(*Element)
			if !ok {
				return nil, fmt.Errorf("invalid index: %d (arg %d), value not a sub-list", k, i)
			}
			v = subTree
		}
	}
	return v.Value, nil
}

func (e *Element) GetFrom(k int) (*Element, error) {
	p := e
	for i := 0; i < k; i++ {
		if p.Next == nil {
			return nil, fmt.Errorf("invalid index: %d, last is %d", k, i)
		}
		p = p.Next
	}
	return p, nil
}

func (e *Element) GetList(kss ...[]int) (*Element, error) {
	res := make([]interface{}, 0, len(kss))
	for _, ks := range kss {
		r, err := e.Get(ks...)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return New(res...), nil
}

func Getter(ks ...int) func(*Element) (interface{}, error) {
	return func(e *Element) (interface{}, error) {
		return e.Get(ks...)
	}
}

func GenericGetter(ks ...int) func(interface{}) (interface{}, error) {
	return func(e interface{}) (interface{}, error) {
		return e.(*Element).Get(ks...)
	}
}

func GetterFrom(k int) func(*Element) (*Element, error) {
	return func(e *Element) (*Element, error) {
		return e.GetFrom(k)
	}
}

func GenericGetterFrom(k int) func(interface{}) (interface{}, error) {
	return func(e interface{}) (interface{}, error) {
		return e.(*Element).GetFrom(k)
	}
}

func ListGetter(kss ...[]int) func(*Element) (*Element, error) {
	return func(e *Element) (*Element, error) {
		return e.GetList(kss...)
	}
}

func GenericListGetter(kss ...[]int) func(interface{}) (interface{}, error) {
	return func(e interface{}) (interface{}, error) {
		return e.(*Element).GetList(kss...)
	}
}

func Cons(value interface{}, next *Element) *Element {
	return &Element{Value: value, Next: next}
}

func ConsGetter(ks1 []int, ks2 []int) func(*Element) (*Element, error) {
	return func(e *Element) (*Element, error) {
		value, err := e.Get(ks1...)
		if err != nil {
			return nil, err
		}
		next, err := e.Get(ks2...)
		if err != nil {
			return nil, err
		}
		return Cons(value, next.(*Element)), nil
	}
}

func GenericConsGetter(ks1 []int, ks2 []int) func(interface{}) (interface{}, error) {
	return func(e interface{}) (interface{}, error) {
		return ConsGetter(ks1, ks2)(e.(*Element))
	}
}
