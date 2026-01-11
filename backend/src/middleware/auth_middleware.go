package middleware

import "github.com/labstack/echo/v4"

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("auth_user_id", "11111111-1111-1111-1111-111111111111")
			return next(c)
		}
	}
}
