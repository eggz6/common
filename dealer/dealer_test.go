package dealer

import (
	"testing"
)

func Test_Shuffle(t *testing.T) {
	d := newDealer52()
	res := d.Shuffle()
	if len(res) != 52 {
		t.Fatal("shuffle err")
	}

	tmp := make(map[int]bool)

	for _, i := range res {
		_, ok := tmp[i]
		if ok {
			t.Fatal("shuffle multiple err")
		}

		tmp[i] = true
	}
}
