package structl

import "testing"

type x string

func (x) M() {}

type y struct {
	z byte
}

type A struct {
	x
}

func (A) y(int) bool {
	return false
}

type B struct {
	y
}

func (B) x(string) {}

func TestAmbiguousEmbedStruct(t *testing.T) {
	var v struct {
		A
		B
	}
	// _ = v.x // error: 模棱两可的v.x
	// _ = v.y // error: 模棱两可的v.y
	_ = v.M // ok. <=> v.A.x.M
	_ = v.z // ok. <=> v.B.y.z
}
