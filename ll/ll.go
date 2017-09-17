/*
Package ll extends container/list to implement the simplest possible double-linked-list with useful methods for accessing deelply nested list (aka "trees"), and decent string representations.

You'll probably need this if you worked with highly mutable trees.

It's `GetIn` and `TryGetIn` are inspired by Clojure's `get-in` but adapted to be more Go-idiomatic. Achieving the funcitonality of methods like `head`, `rest`, `tail`, `car`, `cadr` etc. is possible by sinply using `l.GetIn(0)`, `l.GetIn(1)`, `l.GetIn(0, 1)`, no need to bloat this with extra methods.
*/
package ll

import (
	"container/list"
	"fmt"
	"strings"
)

// List is a duble-linked list with some handy extra methods for workin
// with nested lists of lists ("trees")
type List struct {
	*list.List
}

// New creates a super powered double-linked list
func New(vals ...interface{}) *List {
	head := list.New()
	for _, v := range vals {
		head.PushBack(v)
	}
	return &List{head}
}

// String gives an S-expression style (parens and values) string representation
func (l *List) String() string {
	words := make([]string, 0, 3*l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		words = append(words, fmt.Sprint(e.Value))
	}
	return "(" + strings.Join(words, " ") + ")"
}

// Dump returns a more debugging-friendly string representation
func (l *List) Dump() string {
	words := make([]string, 0, 3*l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		// words = append(words, fmt.Sprint(e.Value))
		switch v := e.Value.(type) {
		case *List:
			words = append(words, v.Dump())
		default:
			words = append(words, fmt.Sprintf("%v::%T", v, v))
		}
	}
	return "[" + strings.Join(words, " ") + "]"
}

// GetIn fetches an element e[ks[0]]...[ks[n]] from a (possible nested) list
// IMPORTANT: this panics when there is no such position in the list
func (l *List) GetIn(ks ...int) interface{} {
	v := l.Front()
	iMax := len(ks) - 1
	for i, k := range ks {
		for j := 0; j < k; j++ {
			v = v.Next()
		}
		if i < iMax {
			v = v.Value.(*List).Front()
		}
	}
	return v.Value
}

// TryGetIn fetches an element e[ks[0]]...[ks[n]] from a (possible nested) list
// or returns defaultVal if no such position exists in the list
func (l *List) TryGetIn(defaultVal interface{}, ks ...int) interface{} {
	v := l.Front()
	iMax := len(ks) - 1
	for i, k := range ks {
		for j := 0; j < k; j++ {
			if v.Next() == nil {
				return defaultVal
			}
			v = v.Next()
		}
		if i < iMax {
			subTree, ok := v.Value.(*List)
			if !ok {
				return defaultVal
			}
			v = subTree.Front()
		}
	}
	return v.Value
}
