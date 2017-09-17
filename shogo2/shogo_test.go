package shogo2

import "github.com/NeuronQ/shoggoth/sll"

var E = sll.New

func makeResult(parserName string, values ...interface{}) interface{} {
	return sll.New(append([]interface{}{parserName}, values...)...)
}
