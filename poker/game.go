package poker

import (
	"context"
	"time"
)

type cmdable interface {
	exec(ctx context.Context, g *game, d time.Duration)
}

type runnable func(ctx context.Context, g *game, d time.Duration)

func (r runnable) exec(ctx context.Context, g *game, d time.Duration) {
	r(ctx, g, d)
}

func idle(ctx context.Context, g *game, d time.Duration) {
}

func prepare(ctx context.Context, g *game, d time.Duration) {
}

func dealing(ctx context.Context, g *game, d time.Duration) {
}

func inHand(ctx context.Context, g *game, d time.Duration) {
}

func fold(ctx context.Context, g *game, d time.Duration) {
}

func turn(ctx context.Context, g *game, d time.Duration) {
}

func river(ctx context.Context, g *game, d time.Duration) {
}

func showhand(ctx context.Context, g *game, d time.Duration) {
}

func settlement(ctx context.Context, g *game, d time.Duration) {
}

func pause(ctx context.Context, g *game, d time.Duration) {
}

func cmdIdle() cmdable {
	return runnable(idle)
}

type game struct {
	runnable
	ID    string
	state GameState
	table *Table
	cmd   chan Cmder
	rend  chan time.Duration
}

func newGame(id string, table *Table) *game {
	res := &game{
		ID:       id,
		state:    Idle,
		runnable: idle,
		table:    table,
		cmd:      make(chan Cmder, 0),
		rend:     make(chan time.Duration, 0),
	}

	return res
}

func (g *game) Start() {
	go func(g *game) {
		for {
			render := false
			select {
			case c, _ := <-g.cmd:
				c.Exec(context.Background(), g)
				// todo broadcast event res
			case r, _ := <-g.rend:
				render = true
				g.runnable(context.Background(), g, r)
				// todo broadcast event res / game info
			}

			if !render {
				select {
				case r, _ := <-g.rend:
					render = true
					g.runnable(context.Background(), g, r)
					// todo broadcast event res / game info
				default:
				}
			}
		}
	}(g)
}

func (g *game) Render(d time.Duration) {
	select {
	case g.rend <- d:
	default:
	}
}
