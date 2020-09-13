package errx

import (
	"fmt"
	"log"
	"testing"
)

func Test_String(t *testing.T) {
	log.Println()
	exx := giveErr()
	exx1 := E(exx)

	t.Log("exx", exx1)
	t.Log("Extract", Extract(exx1))
}

func giveErr() error {
	return E(fmt.Errorf("new error"))
}
