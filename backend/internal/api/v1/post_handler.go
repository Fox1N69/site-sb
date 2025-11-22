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
	GetPostByID(c *gin.Context)
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

	var post models.Post
	if err := c.ShouldBind(&post); err != nil {
		response.New(c).Error(400, err)
		return
	}

	post.ID = uint(postID)

	if err := h.service.UpdatePost(&post); err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, gin.H{
		"message": "post update success",
		"data":    post,
	})
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

func (h *postHandler) GetPostByID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.New(c).Error(400, err)
		return
	}

	post, err := h.service.GetPostByID(uint(postID))
	if err != nil {
		response.New(c).Error(501, err)
		return
	}

	c.JSON(200, post)
}
