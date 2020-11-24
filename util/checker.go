package util

import "fmt"

func Check(ok bool, msg string) func() (bool, string) {
	return func() (bool, string) {
		return ok, msg
	}
}

func Checker(targets ...func() (bool, string)) error {
	var err error

	for _, t := range targets {
		ok, msg := t()
		if !ok {
			err = fmt.Errorf("invalid params. %v", msg)
			break
		}
	}

	return err
}
