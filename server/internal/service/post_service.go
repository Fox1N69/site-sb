package service

import (
	"site-sb/internal/models"
	"site-sb/internal/repo"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetPosts() ([]models.Post, error)
	UpdatePost(postID uint) error
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

func (s *postService) UpdatePost(postID uint) error {
	return s.repo.Update(postID)
}

func (s *postService) DeletePost(postID uint) error {
	return s.repo.Delete(postID)
}
