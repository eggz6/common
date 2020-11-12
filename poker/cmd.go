package poker

import (
	"context"
	"sync/atomic"
	"time"
)

type Cmder interface {
	Exec(ctx context.Context, g *game) error
	Result() Cmdres
}

type Cmdres interface {
	Code() Cmdcode
	Error() error
	Res() interface{}
}

type Gamer interface {
	DoSync(ctx context.Context, c Cmder) Cmdres
	RenderSync(ctx context.Context, d time.Duration)
}

type Result func() Cmdres

func (g *game) DoSync(ctx context.Context, c Cmder) Cmdres {
	g.cmd <- c

	res := c.Result()
	return res
}

func (g *game) RenderSync(ctx context.Context, d time.Duration) {
	g.rend <- d

	return
}

type res struct {
	ch       chan bool
	executed uint32
	data     interface{}
	err      error
	code     Cmdcode
}

func newres() *res {
	return &res{ch: make(chan bool, 1)}
}

func (r *res) tryExecute() bool {
	return atomic.CompareAndSwapUint32(&r.executed, 0, 1)
}

func (r *res) decorator(ctx context.Context, g *game, handle func(ctx context.Context, g *game) error) error {
	if handle == nil {
		return InvalidResHandle
	}

	if !r.tryExecute() {
		return HasBeenExecuted
	}

	defer func() {
		r.ch <- true
		close(r.ch)
	}()

	return handle(ctx, g)
}

func (r *res) Executed() bool {
	abort := atomic.LoadUint32(&r.executed)

	return abort == 1
}

func (r *res) Result() Cmdres {
	<-r.ch

	return r
}

func (r *res) reset() {

}

func (r *res) Error() error {
	if r.err != nil {
		return r.err
	}

	return r.code
}

func (r *res) Res() interface{} {
	return r.data
}

func (r *res) Code() Cmdcode {
	return r.code
}

type Cmdcode uint32

func (e Cmdcode) Error() string {
	res, ok := codeMsg[e]
	if !ok {
		return "unknown"
	}

	return res
}

const (
	success Cmdcode = iota
	invalidSeatIndex
	seatInUsed
	userNoInTable
	userNoInSeat
	userHasInSeat
)

const stand uint32 = 12

var (
	codeMsg map[Cmdcode]string = map[Cmdcode]string{
		success:          "success",
		invalidSeatIndex: "invalid seat index",
		seatInUsed:       "seat in used",
		userNoInTable:    "user no in the table",
		userNoInSeat:     "user no in the seat",
		userHasInSeat:    "user has been in the seat",
	}
)

type CmderErr string

const (
	HasBeenExecuted  = CmderErr("the cmd has been executed")
	InvalidResHandle = CmderErr("invalid res handle")
)

func (e CmderErr) Error() string {
	return string(e)
}
