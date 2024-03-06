package handler

import (
	"rest/database"
	"rest/models"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	//`"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SingUp(c echo.Context) error {
	var data models.Users

	if err := c.Bind(&data); err != nil {
		c.JSON(400, err)
	}

	database.DB.Create(&data)

	return c.JSON(200, data)
}

func (h *Handler) SingIn(c echo.Context) error {
	var data map[string]string

	log.Fatal(c.Bind(&data))

	return nil
}
