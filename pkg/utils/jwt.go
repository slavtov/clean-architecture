package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/slavken/go-clean-architecture/pkg/logger"
)

type JWTConfig struct {
	JWTSecret        string
	JWTRefreshSecret string
	AtExpires        int
	RtExpires        int
}

type TokenDetails struct {
	AtID         uuid.UUID
	RtID         uuid.UUID
	AtExpires    int64
	RtExpires    int64
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(cfg *JWTConfig, id uuid.UUID) (*TokenDetails, error) {
	atID := uuid.New()
	rtID := uuid.New()

	atExpires := getExp(cfg.AtExpires)
	rtExpires := getExp(cfg.RtExpires)

	accessToken, err := createToken(atID, id, atExpires, cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := createToken(rtID, id, rtExpires, cfg.JWTRefreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		AtID:         atID,
		RtID:         rtID,
		AtExpires:    atExpires,
		RtExpires:    rtExpires,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func getExp(seconds int) int64 {
	return time.Now().Add(time.Second * time.Duration(seconds)).Unix()
}

func createToken(
	id uuid.UUID,
	userID uuid.UUID,
	exp int64,
	secret string,
) (string, error) {
	claims := Claims{
		id.String(),
		userID.String(),
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyRefreshToken(
	c echo.Context,
	secret string,
	log logger.Logger,
) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		log.Errorf("c.Cookie: %v", err)
		return err
	}

	if err = ValidateToken(
		c,
		"refresh",
		refreshCookie.Value,
		secret,
	); err != nil {
		log.Errorf("validateToken: %v", err)
		return err
	}

	return nil
}

func ValidateToken(
	c echo.Context,
	tokenName string,
	tokenString string,
	secret string,
) error {
	errString := fmt.Sprintf("invalid %s token", tokenName)

	if tokenString == "" {
		return errors.New(errString)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New(errString)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)
		if !ok {
			return errors.New("invalid jwt claims")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return errors.New("invalid jwt claims")
		}

		tokenUuid, err := uuid.Parse(id)
		if err != nil {
			return err
		}

		userUuid, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		c.Set(fmt.Sprintf("%s_id", tokenName), tokenUuid)
		c.Set("user_id", userUuid)
	}

	return nil
}
