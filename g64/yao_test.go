package g64

import "testing"

func Test_Yao(t *testing.T) {
	g, _ := newGua3(62)

	for i := 5; i >= 0; i-- {
		y := g.Yao(Idx(i))

		if y == 1 {
			t.Logf("---")
			continue
		}

		t.Logf("- -")
	}

	t.Logf("inner=%v, outer=%v", g.Inner(), g.Outer())
}
