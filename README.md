# Shoggoth - simple tools for building complex parsers

![](/logo.jpg =250x250)

## STATUS:

**Work in progress!** Don't use it! Don't even ask me to explain, pretty pease!

## What is this?

Shoggoth is a collection of tools for building advanced parsers, even if you have
**absolutely zero theoretical knowlede** using a techniques called "parser
combinators" which basically means "let's make parsing fuctions, then functions
that return interesting and useful combinations of parsing functions and so on".

## Why?

Because Parser Combinators are the coolest way to write infinitely flexible parsers
(like "parse C++ with JSON sprinkled as comments that changes semantic" or
"parse this partially defined wrong grammar by Friday"), but all/most PC tools are
written by academics who love nothing more than to jerk off with fancy sounding
concepts and sprinkle monds all over their gonads...

Because sometimes you have something that is totally fucked up, ill defined, and
would require a 160+ IQ to write a proper gramar usable by a parser generator for it,
but you still have to write a parser for it, under an unrealistic deadline, and with
the clear requirement that it is not tangle of regexes and has no buffer overflows...

## Install

```
go get https://github.com/NeuronQ/shoggoth
```

## Example

**NOTE:** This is already OUTLDATED, we've moved from struct based to functinal parsers
for simplicity. See `shogo2/` equialents to what's below.

```go
// input string
input := "(+ 123 56 (* 10 (/ 670 30)) 42)"

second := sll.GenericGetter(1)

// "grammar"
var operator Parser
operator = Or("operator", Rune('+'), Rune('-'), Rune('*'), Rune('/'))
operator = Action(operator,
	// (operator (rune +)) -> (operator +)
	sll.GenericListGetter([]int{0}, []int{1, 1}),
)
digit := Or("digit",
	Rune('0'), Rune('1'), Rune('2'), Rune('3'), Rune('4'),
	Rune('5'), Rune('6'), Rune('7'), Rune('8'), Rune('9'))
number := Action(Some("number", digit),
	// (number (digit (rune 4)) (digit (rune 2))) -> (number 42)
	makeNumber)
atom := Action(Or("atom", number, operator), second)

var list Parser
list = Action(
	SpacedSeq("list",
		&RuneParser{r: '('},
		Action(
			SpacedSome("_some_atoms_or_lists_",
				Action(
					Or("_atom_or_list_", atom, Ref(&list)),
					second)),
			// "all but first"
			sll.GenericGetterFrom(1)),
		&RuneParser{r: ')'}),
	// (list ((...) (...) ...)) -> (list (...) (...) ...)
	sll.GenericConsGetter([]int{0}, []int{1}),
)
parser := list

// parsing
pos, r, err := parser.Parse(input, 0, makeResult)

/*
assert.Nil(err)
assert.Equal(31, pos)
assert.Equal("(list (operator +) (number 123) (number 56) (list (operator *) (number 10) (list (operator /) (number 670) (number 30))) (number 42))", r.(*sll.Element).String())
*/
```
