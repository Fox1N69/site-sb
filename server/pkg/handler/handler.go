package handler

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	
}

func (h *Handler) InitRouter() *echo.Echo {
	e := echo.New()

	admin := e.Group("/admin")
	{
		admin.GET("/", h.Admin)
	}

	auth := e.Group("/auth")
	{
		auth.POST("/sing-up", h.SingUp)
		auth.POST("/sing-in", h.SingIn)
	}

	post := e.Group("/posts")
	{
		post.GET("/", h.getAllPosts)
		post.POST("/set", h.setPost)
		post.PUT("/update", h.updatePosts)
		post.DELETE("/delete", h.deletePosts)

	}

	return e
}
