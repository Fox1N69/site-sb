package response

import (
	"net/http"
	"site-sb/internal/models"

	"github.com/gin-gonic/gin"
)

type Wrapper interface {
	Write(code int, message string)
	Error(code int, err error)
	Token(expired string, token string)
}

type wrapper struct {
	c *gin.Context
}

func New(c *gin.Context) Wrapper {
	return &wrapper{c: c}
}

func (w *wrapper) Write(code int, message string) {
	w.c.JSON(code, models.Response{Code: code, Message: message})
}

func (w *wrapper) Error(code int, err error) {
	w.c.JSON(code, models.Response{Code: code, Message: err.Error()})
}

func (w *wrapper) Token(expired string, token string) {
	w.c.JSON(http.StatusOK, models.Token{Expired: expired, Token: token})
}
