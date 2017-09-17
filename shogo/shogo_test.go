package shogo

import (
	"fmt"
	"math"
	"testing"

	"github.com/NeuronQ/shoggoth/sll"
	"github.com/stretchr/testify/assert"
)

func Expect(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("expected %v got %v", expected, got)
	}
}

var E = sll.New

func makeResult(parserName string, values ...interface{}) interface{} {
	return sll.New(append([]interface{}{parserName}, values...)...)
}

// func second(r interface{}) (interface{}, error) {
// 	return r.(*sll.Element).Get(1)
// }

// func rest(r interface{}) (interface{}, error) {
// 	return r.(*sll.Element).Next, nil
// }

/* makeNumber makes a number
 * eg.
 *    (number (digit (rune 4)) (digit (rune 2)))
 *    ->
 *    (number 42)
 */
func makeNumber(r interface{}) (interface{}, error) {
	rList := r.(*sll.Element)
	rLen := rList.Len() - 1
	var n int
	for i, l := 0, rList; l != nil; i, l = i+1, l.Next {
		if i == 0 { // skip first ("number")
			continue
		}
		// ["digit" ["rune" "4"]] -> '4'
		digitStr, err := l.Value.(*sll.Element).Get(1, 1)
		if err != nil {
			return nil, err
		}
		digitByte := digitStr.(string)[0]
		n += int(math.Pow10(rLen-i)) * int(digitByte-'0')
	}
	return E("number", fmt.Sprint(n)), nil
}

func TestShogo(t *testing.T) {
	assert := assert.New(t)

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

	assert.Nil(err)
	assert.Equal(31, pos)
	assert.Equal("(list (operator +) (number 123) (number 56) (list (operator *) (number 10) (list (operator /) (number 670) (number 30))) (number 42))", r.(*sll.Element).String())
}
