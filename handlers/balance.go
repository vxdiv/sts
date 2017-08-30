package handlers

import (
	"net/http"

	"sts/db"
)

type PlayerBalance struct {
	PlayerID string  `json:"playerId"`
	Balance  float64 `json:"balance"`
}

// Player balance
//
// Example: GET /balance?playerId=P1
// Example response: {"playerId": "P1", "balance": 456.00}
func Balance(ctx AppContext) error {
	id := ctx.QueryParam("playerId")
	player, err := db.NewPlayerRepo(ctx.Session).FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)

	}

	balance := PlayerBalance{
		PlayerID: player.ID,
		Balance:  float64(player.Balance) / 100,
	}

	return ctx.JSON(http.StatusOK, balance)
}
