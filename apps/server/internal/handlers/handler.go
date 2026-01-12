package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	db "tourbackend/internal/database/gen"

	"github.com/labstack/echo/v4"
)

// The Handler is a struct which methods handle the individual endpoints
// the handler has a helper method called newReqCtx which adds creates a context for each request
// this context allows interaction with the database and has a few helper methods like Error and JSONMsg
type Handler struct {
	queries    *db.Queries
	IsDeployed bool
}

func NewHandler(queries *db.Queries, isDeployed bool) *Handler {
	return &Handler{queries: queries, IsDeployed: isDeployed}
}

type RequestCtx struct {
	Ctx     context.Context
	Echo    echo.Context
	Queries *db.Queries

	User *db.User
}

func (h *Handler) NewReqCtx(c echo.Context) *RequestCtx {

	r := &RequestCtx{
		Ctx:     c.Request().Context(),
		Echo:    c,
		Queries: h.queries,
	}

	if user, ok := c.Get("user").(*db.User); ok {
		r.User = user
	}

	return r
}

// Helpers

func (r *RequestCtx) Error(code int, msg string) error {
	// perhaps we would like to log the error here in the future?
	return r.Echo.JSON(code, map[string]string{"message": msg})
}

// returns a classic 500 Internal Server Error and logs the error
func (r *RequestCtx) ServerError(err error) error {
	slog.ErrorContext(
		r.Ctx,
		"internal server error",
		"error", err,
	)

	fmt.Println("Unexpected error", err)

	return r.Echo.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
}

func (r *RequestCtx) JSONMsg(code int, msg string) error {
	return r.Echo.JSON(code, map[string]string{"message": msg})
}
