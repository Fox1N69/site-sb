package v1

import (
	"site-sb/common/http/response"
	"site-sb/internal/models"
	"site-sb/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler interface {
	Create(c *gin.Context)
}

type postHandler struct {
	service service.PostService
}

func NewPostHandler(postService service.PostService) PostHandler {
	return &postHandler{service: postService}
}

func (h *postHandler) Create(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBind(post); err != nil {
		response.New(c).Error(400, err)
		return
	}

	if err := h.service.CreatePost(&post); err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, gin.H{
		"message": "post create success",
		"data":    post,
	})
}
