package util

import (
	"strings"
	"testing"
)

func Test_Checker(t *testing.T) {
	a := "1"
	err := Checker(Check(a != "2", "msg error"))

	if err != nil {
		t.Fatalf("should no err here")
	}

	err = Checker(Check(a == "2", "msg error"))
	if err == nil {
		t.Fatalf("need err here")
	}

	if !strings.Contains(err.Error(), "msg error") {
		t.Fatalf("msg failed")
	}
}
