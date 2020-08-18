package rsp

import (
	"github.com/eggz6/common/entity/ret"
)

type M map[string]interface{}

type Result interface {
	Stat() int
	String() string
}

func Rsp(res Result, data interface{}) M {
	return M{
		"code": res.Stat(),
		"data": data,
		"msg":  res.String(),
	}
}

func Success(data interface{}) M {
	res := ret.Success
	return M{
		"code": res.Stat(),
		"msg":  res.String(),
		"data": data,
	}
}
