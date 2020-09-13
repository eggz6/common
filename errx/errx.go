package errx

import (
	"fmt"
	"runtime"
)

type exx struct {
	Err   error
	Local string
}

func (this *exx) Error() string {
	e, ok := this.Err.(*exx)
	if !ok {
		return fmt.Sprintf("%v [%v]", this.Err.Error(), this.Local)
	} else {
		return fmt.Sprintf("%v [%v]", e.Error(), this.Local)
	}
}

func E(err error) error {
	_, f, line, ok := runtime.Caller(1)

	local := ""
	if ok {
		local = fmt.Sprintf("%v %v", f, line)
	}

	res := &exx{Err: err, Local: local}

	return res
}

func Extract(err error) error {
	e, ok := err.(*exx)
	if !ok {
		return err
	}

	return Extract(e.Err)
}
