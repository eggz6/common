package poker

import (
	"context"
	"fmt"
)

// cmdStandUp struct implement Cmder
type cmdStandUp struct {
	res
	userID string
}

func NewCmdStandUp(userID string) Cmder {
	res := newres()
	return &cmdSit{
		res:    *res,
		userID: userID,
	}
}

func (c *cmdStandUp) Exec(ctx context.Context, g *game) error {
	handle := func(ctx context.Context, g *game) error {
		t := g.table
		p, ok := t.players[c.userID]
		if !ok {
			c.code = userNoInTable
			c.err = fmt.Errorf("user=%v no in the table", c.userID)

			return nil
		}

		if p.seat > t.Size {
			c.code = invalidSeatIndex

			return nil
		}

		seat := t.seats[p.seat]
		if seat.playerID != p.playerID {
			c.code = userNoInSeat

			return nil
		}

		seat.playerID = ""
		p.seat = stand

		c.data = seat.id

		return nil
	}

	return c.decorator(ctx, g, handle)
}

func (c *cmdStandUp) Result() Cmdres {
	return c
}
