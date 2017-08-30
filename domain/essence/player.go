package essence

import "errors"

type Player struct {
	ID      string  `json:"playerId"`
	Balance int64   `json:"balance"`
}

var (
	ErrPlayerAlreadyExist = errors.New("Player is already exists")
	ErrPlayerNotFound     = errors.New("Player not found")
)

type PlayerRepository interface {
	Create(player *Player) error
	FindByID(id string) (*Player, error)
	Find(idList []string) ([]*Player, error)
	UpdateBalance(player *Player, amount int64) error
}
