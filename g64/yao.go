package g64

import "fmt"

type Idx uint8

const (
	Chu Idx = iota
	Second
	Third
	Fourth
	Fifth
	Top
)

type Yao uint8

const (
	Jiu Yao = 1
	Liu Yao = 0
)

type Gua6 uint8
type Gua3 uint8

func newGua3(id uint8) (Gua6, error) {
	if id > 63 {
		return 0, fmt.Errorf("invalid id")
	}

	return Gua6(id), nil
}

func (g Gua6) Yao(idx Idx) Yao {
	g8 := uint8(g)

	return Yao((g8 >> idx) & uint8(1))
}

func (g Gua6) Inner() Gua3 {
	return Gua3(uint8(g) & uint8(8))
}

func (g Gua6) Outer() Gua3 {
	return Gua3(uint8(g) >> 3)
}
