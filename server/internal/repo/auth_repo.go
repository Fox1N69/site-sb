package repo

import (
	"site-sb/internal/models"

	"gorm.io/gorm"
)

type AuthRepo interface {
	CheckUsername(username string) bool
	Register(user *models.User) error
	Login(username string) (string, error)
	CheckID(id int) bool
	Delete(id int) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveRecoverToken(user *models.User, token string) error
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *authRepo) CheckUsername(username string) bool {
	var count int64
	if err := r.db.Table("users").Where("username = ?", username).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *authRepo) Register(user *models.User) error {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepo) Login(username string) (string, error) {
	var user models.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (r *authRepo) CheckID(id int) bool {
	var count int64
	if err := r.db.Table("users").Where("id = ?", id).Count(&count).Error; err != nil {
		return false
	}

	if count < 1 {
		return false
	}

	return true
}

func (r *authRepo) Delete(id int) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) SaveRecoverToken(user *models.User, token string) error {
	user.PasswordResetToken = token
	if err := r.db.Save(user).Error; err != nil {
		return nil
	}
	return nil
}
