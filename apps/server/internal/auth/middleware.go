package auth

import (
	"fmt"
	"net/http"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

	"github.com/labstack/echo/v4"
)

//* this file includes middleware - ei code that is run before the actual code for an endpoint.

// Checks if the request includes auth token,
// if it does it validates the token and if the token is valid
// it retrieves from the db data about the user and adds them to the echo context
func AuthMiddleware(queries *db.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			fmt.Println("in auth middleware")

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
				fmt.Println("invalid token")
				return next(c)
			}

			c.Set("user", user)

			fmt.Println("set user")

			return next(c)
		}
	}
}

func AdminRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			fmt.Println("in admin required")

			user, ok := c.Get("user").(*handlers.User)

			if !ok || user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "authentication required",
				})
			}

			if !user.IsAdmin {
				return c.JSON(http.StatusForbidden, map[string]string{
					"message": "admin access required",
				})
			}

			return next(c)
		}
	}
}
