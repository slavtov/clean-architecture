package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slavken/go-clean-architecture/internal/config"
	"github.com/slavken/go-clean-architecture/pkg/logger"
)

func ClearCookies(
	cfg *config.Config,
	log logger.Logger,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetCookie(&http.Cookie{
				Name:   "access_token",
				Path:   "/",
				MaxAge: -1,
			})

			c.SetCookie(&http.Cookie{
				Name:   "refresh_token",
				Path:   "/",
				MaxAge: -1,
			})

			return next(c)
		}
	}
}
