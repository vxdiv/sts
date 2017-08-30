package handlers

import (
	"net/http"

	"sts/db"
)

// Reset DB.
//
// Example: GET /reset Should reset DB to initial state
func Reset(ctx AppContext) error {
	if err := db.Drop(ctx.Session); err != nil {
		return ctx.JSON(http.StatusInternalServerError, Error(err))
	}

	return ctx.String(http.StatusOK, "")
}
