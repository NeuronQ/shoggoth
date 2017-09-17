package shogo2

import "unicode/utf8"

func Rune(r rune, resultName string, name string) Parser {
	return func(input string, pos int, makeResult ASTreeMaker) (newPos int, result interface{}, rerr error) {
		nextRune, w := utf8.DecodeRuneInString(input[pos:])
		if nextRune != r {
			return pos, nil, ParsingError{
				parserName:    name,
				lastParsedPos: pos,
				errPos:        pos + w,
				expected:      string(r),
				found:         input[pos : pos+w],
			}
		}
		if resultName != "" {
			result = makeResult(resultName, string(r))
		}
		return pos + w, result, nil
	}
}
