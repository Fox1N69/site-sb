package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"site-sb/common/http/response"
	"site-sb/common/util/token"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CORS() gin.HandlerFunc
	AUTH() gin.HandlerFunc
	Role(requiredRole string) gin.HandlerFunc
}

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{secretKey: secretKey}
}

func (m *middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *middleware) AUTH() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authBearer := strings.Split(authHeader, " ")

		if len(authBearer) != 2 {
			response.New(c).Error(http.StatusUnauthorized, errors.New("invalid authorization token"))
			c.Abort()
			return
		}

		tokenString := authBearer[1]
		token, err := token.NewToken(m.secretKey).ValidateToken(tokenString)
		if err != nil || !token.Valid {
			response.New(c).Error(http.StatusUnauthorized, fmt.Errorf("invalid authorization token: %v", err))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *middleware) Role(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authBearer := strings.Split(authHeader, " ")

		if len(authBearer) != 2 {
			response.New(c).Error(http.StatusUnauthorized, errors.New("invalid authorization token"))
			c.Abort()
			return
		}

		tokenString := authBearer[1]
		token, err := token.NewToken(m.secretKey).ValidateToken(tokenString)
		if err != nil || !token.Valid {
			response.New(c).Error(http.StatusUnauthorized, fmt.Errorf("invalid authorization token: %v", err))
			c.Abort()
			return
		}

		// Получение роли пользователя из токена
		userRole, ok := token.Claims.(jwt.MapClaims)["role"].(string)
		if !ok {
			response.New(c).Error(http.StatusUnauthorized, errors.New("unable to extract user role from token"))
			c.Abort()
			return
		}

		// Проверка наличия требуемой роли у пользователя
		if userRole != requiredRole {
			response.New(c).Error(http.StatusForbidden, errors.New("access denied: insufficient role"))
			c.Abort()
			return
		}

		c.Next()
	}
}
