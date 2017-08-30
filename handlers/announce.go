package handlers

import (
	"net/http"

	"sts/db"
	"sts/domain/services"
)

type Announce struct {
	TournamentID uint `json:"tournamentId" form:"tournamentId" query:"tournamentId"`
	Deposit      uint `json:"deposit" form:"deposit" query:"deposit"`
}

// Announce tournament specifying the entry deposit
//
// Example: GET /announceTournament?tournamentId=1&deposit=1000
func AnnounceTournament(ctx AppContext) error {
	form := new(Announce)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	service := services.NewTournamentService(
		db.NewTournamentRepo(ctx.Session),
		db.NewPlayerRepo(ctx.Session),
	)

	tournament, err := service.Announce(form.TournamentID, int64(form.Deposit))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	return ctx.JSON(http.StatusOK, tournament)
}
