package v1

import (
	"site-sb/common/http/response"
	"site-sb/internal/models"
	"site-sb/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type postHandler struct {
	service service.PostService
}

func NewPostHandler(postService service.PostService) PostHandler {
	return &postHandler{service: postService}
}

func (h *postHandler) Create(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBind(&post); err != nil {
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

func (h *postHandler) GetAll(c *gin.Context) {
	post, err := h.service.GetPosts()
	if err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, post)
}

func (h *postHandler) Update(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.New(c).Error(400, err)
		return
	}

	if err := h.service.UpdatePost(uint(postID)); err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, "post update seccess")
}

func (h *postHandler) Delete(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.New(c).Error(400, err)
		return
	}

	if err := h.service.DeletePost(uint(postID)); err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, "post delete success")
}
