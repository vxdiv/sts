package handlers

import (
	"net/http"

	"sts/db"
	"sts/domain/services"
)

type Join struct {
	TournamentID uint     `json:"tournamentId" form:"tournamentId" query:"tournamentId"`
	PlayerID     string   `json:"playerId" form:"playerId" query:"playerId"`
	Backers      []string `json:"backerId" form:"backerId" query:"backerId"`
}

// Join player into a tournament and is he backed by a set of backers
//
// Example: GET /joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3
func JoinTournament(ctx AppContext) error {
	form := new(Join)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	service := services.NewTournamentService(
		db.NewTournamentRepo(ctx.Session),
		db.NewPlayerRepo(ctx.Session),
	)

	if err := service.Join(form.TournamentID, form.PlayerID, form.Backers...); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	return ctx.JSON(http.StatusOK, form)
}
