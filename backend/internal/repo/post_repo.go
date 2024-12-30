package repo

import (
	"site-sb/internal/models"

	"gorm.io/gorm"
)

type PostRepo interface {
	Create(post *models.Post) error
	GetPosts() ([]models.Post, error)
	PostByID(postID uint) (*models.Post, error)
	Update(postID uint) error
	Delete(postID uint) error
}

type postRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) PostRepo {
	return &postRepo{db: db}
}

func (r *postRepo) Create(post *models.Post) error {
	return r.db.Create(&post).Error
}

func (r *postRepo) GetPosts() ([]models.Post, error) {
	var post []models.Post
	if err := r.db.Find(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepo) Update(postID uint) error {
	var post models.Post
	if err := r.db.Where("id = ?", postID).Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepo) Delete(postID uint) error {
	var post models.Post
	if err := r.db.Where("id = ?", postID).Delete(&post).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepo) PostByID(postID uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
