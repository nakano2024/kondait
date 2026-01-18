package middleware

import (
	"errors"
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
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHead, prefix) {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			authToken := strings.TrimPrefix(authHead, prefix)

			if authToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			output, err := getPrincipalUsecase.Exec(c.Request().Context(), usecase.GetPrincipalInput{
				AuthToken: authToken,
			})

			if err != nil {
				var tokenInvalidErr *usecase.TokenInvalidError
				if errors.As(err, &tokenInvalidErr) {
					return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}

			c.Set(dto.PrincipalContextKeyName, dto.Principal{
				ActorCode: output.ActorCode,
				Scopes:    output.Scopes,
			})
			return next(c)
		}
	}
}
