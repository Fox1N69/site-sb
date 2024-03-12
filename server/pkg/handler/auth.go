package handler

import (
	"encoding/json"
	"net/http"
	"rest/database"
	"rest/pkg/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "secret"

func (h *Handler) SingUp(c echo.Context) error {
	data := new(models.Users)

	if err := c.Bind(&data); err != nil {
		c.JSON(400, err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	user := models.Users{
		Username: data.Username,
		Password: password,
	}

	database.DB.Create(&user)

	return c.JSON(200, user)
}

func (h *Handler) SingIn(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return c.JSON(400, err)
	}

	var user models.Users

	if err := database.DB.Where("Username = ?", data["Username"]).First(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.JSON(http.StatusBadRequest, "incorrect password")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(), //частота обновления токена
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "could not generate token")
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	c.SetCookie(&cookie)

	return c.JSON(200, "success")
}

func (h *Handler) getAllUsers(c echo.Context) error {
	var data []models.Users
	database.DB.Find(&data)

	user, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	jsonString := string(user)

	return c.JSON(200, jsonString)
}

func (h *Handler) deleteUser(c echo.Context) error {
	var user models.Users
	id := c.Param("id")
	if id == "" {
		c.JSON(400, "ID empty")
	}

	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return c.JSON(400, "user not found")
	}

	if err := c.Bind(&user); err != nil {
		log.Fatal(err)
	}

	if err := database.DB.Delete(&user, id).Error; err != nil {
		return c.JSON(400, "user not delete")
	}

	return c.JSON(http.StatusOK, "User delete")
}
