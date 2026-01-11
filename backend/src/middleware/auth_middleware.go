package middleware

import (
	"github.com/labstack/echo/v4"

	"kondait-backend/dto/auth"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("principal", auth.Principal{
				UserCode: "11111111-1111-1111-1111-111111111111",
				Scopes: []string{
					auth.ScopeCookingItemsRead,
				},
			})
			return next(c)
		}
	}
}
