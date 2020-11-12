package poker

import "github.com/eggz6/common/dealer"

type tableErr string

func (e tableErr) Error() string {
	return string(e)
}

const InvalidSeats tableErr = "the size must be 2 < size <= 12"

type PlayerState int
type TableState int
type GameState int

const (
	Idle       GameState = iota
	Prepare              // 初始化
	Dealing              // 发牌
	InHand               // 手牌
	Fold                 //翻牌
	Turn                 //转牌
	River                //河牌
	ShowHand             // 摊牌
	Settlement           //结算
	Pause
)

type Table struct {
	ID      uint32
	Size    uint32
	seats   []*seat
	dealer  dealer.Dealer
	players map[string]*player
}

type bet struct {
	gameID  string
	num     uint32
	players []string
}

type seat struct {
	id       uint32
	playerID string
	Bets     uint32
}

type player struct {
	playerID  string
	tableID   uint32
	seat      uint32
	TotalBets uint32

	state PlayerState
}

func NewTable(size uint32) (*Table, error) {
	return newTable(size)
}

func newTable(size uint32) (*Table, error) {
	if size < 2 || size > 12 {
		return nil, InvalidSeats
	}

	seats := make([]*seat, size)
	for i := uint32(0); i < size; i++ {
		seats[i] = &seat{id: uint32(i)}
	}

	return &Table{
		Size:    size,
		seats:   seats,
		dealer:  dealer.NewDealer52(),
		players: make(map[string]*player, 0),
	}, nil
}
