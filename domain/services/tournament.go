package services

import (
	"errors"

	"sts/domain/essence"

	"github.com/labstack/gommon/log"
)

var (
	ErrNotEnoughPoints                = errors.New("Not enough points")
	ErrInvalidWinnerPlayers           = errors.New("Invalid winner player")
	ErrPrizeAmountIsTooLarge          = errors.New("The prize amount is too large")
	ErrTournamentBalanceLimitExceeded = errors.New("The tournament balance limit is exceeded")
	ErrTournamentAlreadyClose         = errors.New("The tournament is already close")
)

type TournamentService struct {
	tournaments essence.TournamentRepository
	players     essence.PlayerRepository
}

func NewTournamentService(tournaments essence.TournamentRepository, players essence.PlayerRepository) *TournamentService {
	return &TournamentService{
		tournaments: tournaments,
		players:     players,
	}
}

func (ts *TournamentService) Announce(id uint, deposit int64) (*essence.Tournament, error) {
	tournament := &essence.Tournament{
		ID:            id,
		DepositPoints: deposit,
		Balance:       0.0,
		Teams:         make(essence.Teams),
	}

	if err := ts.tournaments.Create(tournament); err != nil {
		return nil, err
	}

	return tournament, nil
}

func (ts *TournamentService) Join(tournamentID uint, playerID string, backerList ...string) error {
	tournament, err := ts.tournaments.FindByID(tournamentID)
	if err != nil {
		return err
	}

	player, err := ts.players.FindByID(playerID)
	if err != nil {
		return err
	}

	backers, err := ts.players.Find(backerList)
	if err != nil {
		return err
	}

	team := ts.createTeam(player, backers...)
	depositAmount := ts.depositAmount(tournament, team)

	if err := ts.validateDepositsPart(team, depositAmount); err != nil {
		return err
	}

	if err := ts.withdrawAmountFromTeam(team, depositAmount); err != nil {
		return err
	}

	team.WinPercent = ts.calcPercentPart(tournament, depositAmount)

	tournament.Balance += depositAmount * team.Count()
	tournament.Teams[player.ID] = team

	if err := ts.tournaments.Update(tournament); err != nil {
		return err
	}

	return nil
}

func (ts *TournamentService) createTeam(player *essence.Player, backers ...*essence.Player) *essence.Team {
	return &essence.Team{
		Player:  player,
		Backers: backers,
	}
}

func (ts *TournamentService) depositAmount(tournament *essence.Tournament, team *essence.Team) int64 {
	points := float64(tournament.DepositPoints) / float64(team.Count())
	return int64(points * 100)
}

func (ts *TournamentService) validateDepositsPart(team *essence.Team, depositAmount int64) error {
	if team.CountBackers() == 0 && team.Player.Balance < depositAmount {
		return ErrNotEnoughPoints
	}

	for _, backer := range team.Backers {
		if backer.Balance < depositAmount {
			log.Errorf("BACKER ID %v not enough balance. Balance %v", backer.ID, backer.Balance)

			return ErrNotEnoughPoints
		}
	}

	return nil
}

func (ts *TournamentService) withdrawAmountFromTeam(team *essence.Team, amount int64) error {
	if err := ts.players.UpdateBalance(team.Player, -amount); err != nil {
		return err
	}

	for _, backer := range team.Backers {
		if err := ts.players.UpdateBalance(backer, -amount); err != nil {
			return err
		}
	}

	return nil
}

func (ts *TournamentService) calcPercentPart(tournament *essence.Tournament, amount int64) int64 {
	return amount / tournament.DepositPoints
}

type Winners map[string]int64

func (w Winners) PrizeSum() int64 {
	sum := int64(0)
	for _, prize := range w {
		sum += prize
	}

	return sum
}

type WinnerResult struct {
	Winners []WinnerPlayer `json:"winners"`
}

type WinnerPlayer struct {
	PlayerID string `json:"playerId"`
	Prize    int64  `json:"prize"`
}

func (ts *TournamentService) ProcessResult(tournamentID uint, winners Winners) (*WinnerResult, error) {
	tournament, err := ts.tournaments.FindByID(tournamentID)
	if err != nil {
		return nil, err
	}

	if tournament.Close {
		return nil, ErrTournamentAlreadyClose
	}

	if tournament.Balance < winners.PrizeSum() {
		return nil, ErrPrizeAmountIsTooLarge
	}

	winnersResult := &WinnerResult{
		Winners: make([]WinnerPlayer, 0, 0),
	}
	for playerId, prize := range winners {
		team, ok := tournament.Teams[playerId]
		if !ok {
			log.Errorf("Player %v not participate in tournament", playerId)

			return nil, ErrInvalidWinnerPlayers
		}

		tournamentCashOut, err := ts.processTeam(team, prize)
		if err != nil {
			return nil, err
		}

		tournament.Balance -= tournamentCashOut
		if tournament.Balance < 0 {
			return nil, ErrTournamentBalanceLimitExceeded
		}

		winnersResult.Winners = append(winnersResult.Winners, WinnerPlayer{PlayerID: playerId, Prize: prize})
	}

	tournament.Close = true
	if err := ts.tournaments.Update(tournament); err != nil {
		return nil, err
	}

	return winnersResult, nil
}

func (ts *TournamentService) processTeam(team *essence.Team, prize int64) (int64, error) {
	partPrize := partPrize(prize, team.WinPercent)
	if err := ts.players.UpdateBalance(team.Player, partPrize); err != nil {
		return 0, err
	}

	for _, backer := range team.Backers {
		if err := ts.players.UpdateBalance(backer, partPrize); err != nil {
			return 0, err
		}
	}

	return partPrize * team.Count(), nil
}

func partPrize(prizePoints, percent int64) int64 {
	return int64(float64(prizePoints) * float64(percent))
}
