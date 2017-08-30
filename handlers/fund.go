package handlers

import (
	"net/http"

	"sts/db"

	"sts/domain/services"
)

type Player struct {
	PlayerID string `json:"playerId" form:"playerId" query:"playerId"`
	Points   int64   `json:"points" form:"points" query:"points"`
}

// Fund player account
//
// Example: GET /fund?playerId=P2&points=300 funds (add to balance) player P2 with 300 points.
// If no player exist should create new player with given amount of points
func Fund(ctx AppContext) error {
	form := new(Player)
	if err := ctx.Bind(form); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	service := services.NewPlayerService(db.NewPlayerRepo(ctx.Session))
	if err := service.IncreasePoints(form.PlayerID, form.Points); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error(err))
	}

	return ctx.String(http.StatusOK, "")
}
