package service

import (
	"site-sb/internal/models"
	"site-sb/internal/repo"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetPosts() ([]models.Post, error)
	GetPostByID(postID uint) (*models.Post, error)
	UpdatePost(post *models.Post) error
	DeletePost(postID uint) error
}

type postService struct {
	repo repo.PostRepo
}

func NewPostService(repository repo.PostRepo) PostService {
	return &postService{repo: repository}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.repo.Create(post)
}

func (s *postService) GetPosts() ([]models.Post, error) {
	return s.repo.GetPosts()
}

func (s *postService) UpdatePost(post *models.Post) error {
	return s.repo.Update(post)
}

func (s *postService) DeletePost(postID uint) error {
	return s.repo.Delete(postID)
}

func (s *postService) GetPostByID(postID uint) (*models.Post, error) {
	return s.repo.PostByID(postID)
}
