package ll

import (
	"testing"
)

func Expect(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Error("expected", expected, "got", got)
	}
}

var E = New

func TestLl_String(t *testing.T) {
	Expect(t, E().String(), "()")
	Expect(t, E(nil).String(), "(<nil>)")
	Expect(t, E(42).String(), "(42)")
	Expect(t, E(11, 22).String(), "(11 22)")
	Expect(t, E("foo", 10).String(), "(foo 10)")
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
	Expect(t, l.String(), "(11 22 ((go [42 34 66]) -3) 33 ok)")
}

func TestSll_Dump(t *testing.T) {
	Expect(t, E().Dump(), "[]")
	Expect(t, E(nil).Dump(), "[<nil>::<nil>]")
	Expect(t, E(42).Dump(), "[42::int]")
	Expect(t, E(11, 22).Dump(), "[11::int 22::int]")
	Expect(t, E("foo", 10).Dump(), "[foo::string 10::int]")
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
	Expect(t, l.Dump(), "[11::int 22::int [[go::string [42 34 66]::[]int] -3::int] 33::int ok::string]")
}

func TestSll_GetIn(t *testing.T) {
	Expect(t, E(nil).GetIn(0), nil)
	Expect(t, E(42).GetIn(0), 42)
	Expect(t, E(11, "howdy").GetIn(1), "howdy")
	Expect(t, E(11, "howdy", -3.14).GetIn(2), -3.14)

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
	Expect(t, l.GetIn(0), 11)
	Expect(t, l.GetIn(4), "ok")
	Expect(t, l.GetIn(2, 1), -3)
	Expect(t, l.GetIn(2, 0, 0), "go")
	Expect(t, l.GetIn(2, 0, 1).([]int)[0], 42)
}

func TestSll_TryGetIn(t *testing.T) {
	Expect(t, E(nil).TryGetIn(-1, 0), nil)
	Expect(t, E(nil).TryGetIn(-1, 1), -1)
	Expect(t, E(42).TryGetIn(-1, 0), 42)
	Expect(t, E(42).TryGetIn(-1, 1), -1)
	Expect(t, E(42).TryGetIn(-1, 12), -1)
	Expect(t, E(11, "howdy").TryGetIn(-1, 1), "howdy")
	Expect(t, E(11, "howdy").TryGetIn(-1, 2), -1)
	Expect(t, E(11, "howdy", -3.14).TryGetIn(-1, 2), -3.14)
	Expect(t, E(11, "howdy", -3.14).TryGetIn(-1, 3), -1)

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
	Expect(t, l.TryGetIn(-1, 0), 11)
	Expect(t, l.TryGetIn(-1, 4), "ok")
	Expect(t, l.TryGetIn(-1, 5), -1)
	Expect(t, l.TryGetIn(-1, 2, 1), -3)
	Expect(t, l.TryGetIn(-1, 2, 2), -1)
	Expect(t, l.TryGetIn(-1, 2, 0, 0), "go")
	Expect(t, l.TryGetIn(-1, 2, 0, 2), -1)
	Expect(t, l.TryGetIn(-1, 2, 0, 1).([]int)[0], 42)
	Expect(t, l.TryGetIn(-1, 2, 0, 2, 0), -1)
	Expect(t, l.TryGetIn(-1, 2, 0, 2, 1, 0), -1)
}
