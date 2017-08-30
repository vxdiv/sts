package essence

import "errors"

type Tournament struct {
	ID            uint
	DepositPoints int64
	Balance       int64
	Teams         Teams
	Close         bool
}

type Teams map[string]*Team

type Team struct {
	Player     *Player
	Backers    []*Player
	WinPercent int64
}

func (t Team) Count() int64 {
	return t.CountBackers() + 1
}

func (t Team) CountBackers() int64 {
	return int64(len(t.Backers))
}

var (
	ErrTournamentNotFound     = errors.New("Tournament not found")
	ErrTournamentAlreadyExist = errors.New("Tournament is already exists")
)

type TournamentRepository interface {
	Create(tournament *Tournament) error
	FindByID(id uint) (*Tournament, error)
	Update(tournament *Tournament) error
}
