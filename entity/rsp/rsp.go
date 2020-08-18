package rsp

import (
	"github.com/eggz6/common/entity/ret"
	"github.com/gin-gonic/gin"
)

type Result interface {
	Stat() int
	String() string
}

func Rsp(res Result, data interface{}) gin.H {
	return gin.H{
		"code": res.Stat(),
		"data": data,
		"msg":  res.String(),
	}
}

func Success(data interface{}) gin.H {
	res := ret.Success
	return gin.H{
		"code": res.Stat(),
		"msg":  res.String(),
		"data": data,
	}
}
