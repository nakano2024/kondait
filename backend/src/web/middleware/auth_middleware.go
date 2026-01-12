package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"kondait-backend/application/usecase"
	"kondait-backend/web/dto"
)

func AuthMiddleware(getPrincipalUsecase usecase.IGetPrincipalUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if getPrincipalUsecase == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "auth dependency not set")
			}

			authHead := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHead == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			const Prefix = "Bearer "

			if !strings.HasPrefix(authHead, Prefix) {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			authToken := strings.TrimPrefix(authHead, Prefix)
			output, err := getPrincipalUsecase.Exec(usecase.GetPrincipalInput{
				AuthToken: authToken,
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			c.Set("principal", dto.Principal{
				ActorCode: output.UserCode,
				Scopes:    output.Scopes,
			})
			return next(c)
		}
	}
}
