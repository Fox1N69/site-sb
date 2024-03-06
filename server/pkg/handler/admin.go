package handler

import "github.com/labstack/echo/v4"

func (h *Handler) Admin(c echo.Context) error {
	return c.String(200,"Hello echo")
} 