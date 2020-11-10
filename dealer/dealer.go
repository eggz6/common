package dealer

import (
	"math/rand"
	"time"
)

type Dealer interface {
	Shuffle() []int
}

type dealer52 struct {
}

func NewDealer52() Dealer {
	return newDealer52()
}

func newDealer52() *dealer52 {
	return &dealer52{}
}

func (d *dealer52) Shuffle() []int {
	res := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
		51, 52,
	}

	end := len(res)
	r := d.rand()

	for end > 0 {
		i := r.Intn(end)
		target := end - 1
		tmp := res[target]
		res[target] = res[i]
		res[i] = tmp

		end--
	}

	return res
}

func (d *dealer52) rand() *rand.Rand {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)

	return r
}
