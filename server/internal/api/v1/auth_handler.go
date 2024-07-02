package v1

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"site-sb/common/http/response"
	"site-sb/common/util/regex"
	"site-sb/common/util/token"
	"site-sb/infra"
	"site-sb/internal/models"
	"site-sb/internal/service"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Delete(c *gin.Context)
	GetUsernameByID(c *gin.Context)
}

type authUserHandler struct {
	authService service.AuthService
	infra       infra.Infra
}

func NewAuthHandler(authService service.AuthService, infra infra.Infra) AuthUserHandler {
	return &authUserHandler{
		authService: authService,
		infra:       infra,
	}
}

// GetUserByID retrieves a user by their ID.
func (h *authUserHandler) GetUsernameByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	userID, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *authUserHandler) Register(c *gin.Context) {
	var data models.User
	if err := c.ShouldBind(&data); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
		return
	}

	if err := validation.Validate(data.FIO, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.authService.CheckUsername(data.Username) {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
			return
		}

		data.Password = string(password)
		if err := h.authService.Register(&data); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		response.New(c).Write(http.StatusCreated, "success: user registered")
		return

	}

	response.New(c).Error(http.StatusBadRequest, errors.New("username: already taken"))
}

func (h *authUserHandler) Login(c *gin.Context) {
	var data models.Login
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
		return
	}

	hashedPassword, err := h.authService.Login(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(data.Password)); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("username or password not match"))
		return
	}

	user, err := h.authService.GetUserByUsername(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("username not match"))
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("role", user.Role)
	session.Set("fio", user.FIO)
	session.Save()

	expired, token := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateToken(data.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"fio":      user.FIO,
			"role":     user.Role,
		},
		"token":   token,
		"expired": expired,
	})
}

func (h *authUserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	if h.authService.CheckID(id) {
		if err := h.authService.Delete(id); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		response.New(c).Write(http.StatusOK, "success: user deleted")
		return
	}

	response.New(c).Error(http.StatusBadRequest, errors.New("id: not found"))
}
