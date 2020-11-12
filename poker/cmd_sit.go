package poker

import (
	"context"
	"fmt"
)

// cmdSit struct implement Cmder
type cmdSit struct {
	res
	pos    uint32
	userID string
	bets   uint32
}

func NewCmdSit(pos, bets uint32, userID string) Cmder {
	res := newres()
	return &cmdSit{
		res:    *res,
		pos:    pos,
		bets:   bets,
		userID: userID,
	}
}

func (c *cmdSit) Exec(ctx context.Context, g *game) error {
	handle := func(ctx context.Context, g *game) error {
		t := g.table
		if c.pos >= t.Size {
			c.code = invalidSeatIndex
			c.err = fmt.Errorf("inivalid seat pos=%v", c.pos)

			return nil
		}

		p, ok := t.players[c.userID]
		if ok && p.seat <= stand {
			c.code = userHasInSeat
			c.err = fmt.Errorf("user=%v has seat in pos=%v ", c.userID, p.seat)

			return nil
		}

		seat := t.seats[c.pos]
		if seat.playerID != "" { // todo 自动使用下一个可用座位
			c.code = seatInUsed
			c.err = fmt.Errorf("seat pos=%v has been in used", c.pos)

			return nil
		}

		p = &player{
			playerID:  c.userID,
			tableID:   t.ID,
			seat:      seat.id,
			TotalBets: c.bets,
			state:     0,
		}

		t.players[p.playerID] = p
		seat.playerID = p.playerID
		c.data = seat.id

		return nil
	}

	return c.decorator(ctx, g, handle)
}
