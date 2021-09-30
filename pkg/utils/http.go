package utils

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetCtxID(c echo.Context) uuid.UUID {
	return c.Get("user_id").(uuid.UUID)
}

func GetCtxAccessID(c echo.Context) uuid.UUID {
	return c.Get("access_id").(uuid.UUID)
}

func GetCtxRefreshID(c echo.Context) uuid.UUID {
	return c.Get("refresh_id").(uuid.UUID)
}
