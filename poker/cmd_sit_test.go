package poker

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func Test_Sit(t *testing.T) {
	table, _ := newTable(3)
	g := newGame("1", table)
	g.Start()

	sit := NewCmdSit(1, 100, "user_1")
	res := g.DoSync(context.Background(), sit)
	if res.Code() != success {
		t.Fatalf("sit do sync failed.")
	}

	if fmt.Sprintf("%v", res.Res()) != "1" {
		t.Fatalf("sit do sync sit failed. pos should=%v, actual=%v", 1, res.Res())
	}

	seat := table.seats[1]
	if seat.playerID != "user_1" {
		t.Fatalf("sit do sync sit failed. user_id should=%v, actual=%v", 1, res.Res())
	}

	sit2 := NewCmdSit(1, 100, "user_2")
	res = g.DoSync(context.Background(), sit2)
	if res.Code() != seatInUsed {
		t.Fatalf("sit do sync failed. seat should be in used.")
	}

	sit3 := NewCmdSit(2, 100, "user_1")
	res = g.DoSync(context.Background(), sit3)
	if res.Code() != userHasInSeat {
		t.Fatalf("sit do sync failed. user has been in seat.")
	}
}

func Test_P_Sit(t *testing.T) {
	table, _ := newTable(10)
	g := newGame("1", table)
	g.Start()

	var wg sync.WaitGroup

	t.Run("group", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			pos := uint32(i)
			name := fmt.Sprintf("test_%v", i)
			t.Run(name, func(tt *testing.T) {
				tt.Parallel()
				sit := NewCmdSit(pos, 100, name)
				res := g.DoSync(context.Background(), sit)
				pos32, _ := res.Res().(uint32)
				if pos32 != pos {
					t.Fatalf("parallel do sync sit failed. pos should=%v, actual=%v", pos, res.Res())
				}
				wg.Done()
			})
		}
	})

	wg.Wait()

	for i, seat := range table.seats {
		name := fmt.Sprintf("test_%v", i)
		if seat.playerID != name {
			t.Fatalf("parallel do sync sit failed. usershould=%v, actual=%v", name, seat.playerID)
		}
	}
}
