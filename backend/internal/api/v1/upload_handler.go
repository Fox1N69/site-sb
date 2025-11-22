package v1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"site-sb/common/http/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler interface {
	UploadImages(c *gin.Context)
	DeleteImage(c *gin.Context)
}

type uploadHandler struct {
	uploadPath string
}

func NewUploadHandler(uploadPath string) UploadHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}
	return &uploadHandler{uploadPath: uploadPath}
}

// UploadImages handles multiple file uploads
func (h *uploadHandler) UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.New(c).Error(400, fmt.Errorf("failed to parse multipart form: %v", err))
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		response.New(c).Error(400, fmt.Errorf("no files uploaded"))
		return
	}

	// Limit to 10 images
	if len(files) > 10 {
		response.New(c).Error(400, fmt.Errorf("too many files, maximum is 10"))
		return
	}

	var uploadedUrls []string
	var errors []string

	for _, file := range files {
		// Validate file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			errors = append(errors, fmt.Sprintf("%s: invalid file type, only jpg, jpeg, png, gif, webp allowed", file.Filename))
			continue
		}

		// Validate file size (max 10MB)
		if file.Size > 10*1024*1024 {
			errors = append(errors, fmt.Sprintf("%s: file too large, maximum is 10MB", file.Filename))
			continue
		}

		// Generate unique filename
		uniqueID := uuid.New().String()
		timestamp := time.Now().Format("20060102")
		filename := fmt.Sprintf("%s_%s%s", timestamp, uniqueID, ext)
		filepath := filepath.Join(h.uploadPath, filename)

		// Save file
		src, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to open file", file.Filename))
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to create file", file.Filename))
			continue
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to save file", file.Filename))
			os.Remove(filepath) // Clean up partial file
			continue
		}

		// Return URL path (relative to server)
		url := fmt.Sprintf("/uploads/%s", filename)
		uploadedUrls = append(uploadedUrls, url)
	}

	c.JSON(200, gin.H{
		"message": "upload completed",
		"urls":    uploadedUrls,
		"errors":  errors,
		"count":   len(uploadedUrls),
	})
}

// DeleteImage deletes an uploaded image
func (h *uploadHandler) DeleteImage(c *gin.Context) {
	imageUrl := c.Query("url")
	if imageUrl == "" {
		response.New(c).Error(400, fmt.Errorf("url parameter required"))
		return
	}

	// Extract filename from URL
	filename := filepath.Base(imageUrl)

	// Security check: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		response.New(c).Error(400, fmt.Errorf("invalid filename"))
		return
	}

	filepath := filepath.Join(h.uploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		response.New(c).Error(404, fmt.Errorf("file not found"))
		return
	}

	// Delete file
	if err := os.Remove(filepath); err != nil {
		response.New(c).Error(500, fmt.Errorf("failed to delete file: %v", err))
		return
	}

	c.JSON(200, gin.H{
		"message": "image deleted successfully",
	})
}
