package auth

import (
	db "tourbackend/internal/database/gen"

	"github.com/labstack/echo/v4"
)

//* this file includes middleware - ei code that is run before the actual code for an endpoint.

// Checks if the request includes auth token,
// if it does it validates the token and if the token is valid
// it retrieves from the db data about the user and adds them to the echo context
func AuthMiddleware(queries *db.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ctx := c.Request().Context()

			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return next(c)
			}

			if cookie.Value == "" {
				return next(c)
			}

			user, err := validateToken(cookie.Value, queries, ctx)
			if err != nil {
				// fmt.Println("invalid token")
				return next(c)
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
