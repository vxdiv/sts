package handlers

import (
	"net/http"

	"strconv"

	"sts/db"
	"sts/domain/services"
)

type Result struct {
	TournamentID string     `json:"tournamentId"`
	Winners      []Winner `json:"winners"`
}

type Winner struct {
	PlayerID string `json:"playerId"`
	Prize    int64  `json:"prize"`
}

// Result tournament winners and prizes
//
// Example: POST /resultTournament with body in JSON format
// Example response: {"tournamentId": "1", "winners": [{"playerId": "P1", "prize": 500}]}
func ResultTournament(ctx AppContext) error {
	form := new(Result)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	service := services.NewTournamentService(
		db.NewTournamentRepo(ctx.Session),
		db.NewPlayerRepo(ctx.Session),
	)

	tournamentID, err := strconv.ParseUint(form.TournamentID, 0, 64)
	if err != nil {
		return err
	}

	if err := service.ProcessResult(uint(tournamentID), convertWinners(form.Winners)); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	return ctx.JSON(http.StatusOK, "")
}

func convertWinners(winners []Winner) services.Winners {
	winnerList := make(services.Winners)

	for _, winner := range winners {
		winnerList[winner.PlayerID] = winner.Prize
	}

	return winnerList
}
