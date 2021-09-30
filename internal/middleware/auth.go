package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavken/go-clean-architecture/internal/config"
	"github.com/slavken/go-clean-architecture/internal/domain/usecases"
	"github.com/slavken/go-clean-architecture/pkg/logger"
	"github.com/slavken/go-clean-architecture/pkg/utils"
)

func Auth(
	cfg *config.Config,
	userUseCase usecases.UserUseCase,
	log logger.Logger,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := verifyAccessToken(cfg, c, log); err != nil {
				log.Errorf("verifyAccessToken: %v", err)
				return echo.ErrUnauthorized
			}

			if err := utils.VerifyRefreshToken(
				c,
				cfg.Server.JwtRefreshSecret,
				log,
			); err != nil {
				log.Errorf("verifyRefreshToken: %v", err)
				return echo.ErrUnauthorized
			}

			userID := utils.GetCtxID(c)
			accessID := utils.GetCtxAccessID(c)
			refreshID := utils.GetCtxRefreshID(c)

			if err := verifyRedis(
				userID,
				userUseCase,
				&utils.TokenDetails{
					AtID: accessID,
					RtID: refreshID,
				},
				log,
			); err != nil {
				log.Errorf("verifyRedis: %v", err)
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

func verifyAccessToken(
	cfg *config.Config,
	c echo.Context,
	log logger.Logger,
) error {
	tokenName := "access"
	bearerToken := c.Request().Header.Get("Authorization")
	if bearerToken != "" {
		arr := strings.Split(bearerToken, " ")
		if len(arr) == 2 && arr[0] == "Bearer" {
			token := arr[1]

			if err := utils.ValidateToken(
				c,
				tokenName,
				token,
				cfg.Server.JwtSecret,
			); err != nil {
				log.Errorf("validateToken: %v", err)
				return err
			}

			return nil
		}
	}

	accessCookie, err := c.Cookie("access_token")
	if err != nil {
		log.Errorf("c.Cookie: %v", err)
		return err
	}

	if err = utils.ValidateToken(
		c,
		tokenName,
		accessCookie.Value,
		cfg.Server.JwtSecret,
	); err != nil {
		log.Errorf("validateToken: %v", err)
		return err
	}

	return nil
}

func verifyRedis(
	id uuid.UUID,
	u usecases.UserUseCase,
	td *utils.TokenDetails,
	log logger.Logger,
) error {
	atUserID, err := u.GetToken(id, td.AtID)
	if err != nil {
		log.Errorf("auth.UseCase.GetToken: %v", err)
		return err
	}

	rtUserID, err := u.GetToken(id, td.RtID)
	if err != nil {
		log.Errorf("auth.UseCase.GetToken: %v", err)
		return err
	}

	if atUserID != rtUserID || atUserID != id {
		return echo.ErrForbidden
	}

	return nil
}
