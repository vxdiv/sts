package handlers

import (
	"net/http"

	"sts/db"
	"sts/domain/services"
)

// Take player account
//
// Example: GET /take?playerId=P1&points=300 takes 300 points from player P1 account
func Take(ctx AppContext) error {
	form := new(Player)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	service := services.NewPlayerService(db.NewPlayerRepo(ctx.Session))
	if err := service.IncreasePoints(form.PlayerID, -form.Points); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	return ctx.String(http.StatusOK, "")
}
