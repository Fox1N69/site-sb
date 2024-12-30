package service

import (
	"site-sb/internal/models"
	"site-sb/internal/repo"
)

type AuthService interface {
	CheckUsername(username string) bool
	Register(user *models.User) error
	Login(username string) (string, error)
	CheckID(id int) bool
	Delete(id int) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveRecoveryTokenToDB(user *models.User, token string) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) GetUserByUsername(username string) (*models.User, error) {
	return s.authRepo.GetUserByUsername(username)
}

func (s *authService) GetUserByID(id uint) (*models.User, error) {
	return s.authRepo.GetUserByID(id)
}

func (s *authService) CheckUsername(username string) bool {
	return s.authRepo.CheckUsername(username)
}

func (s *authService) Register(user *models.User) error {
	if err := s.authRepo.Register(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(username string) (string, error) {
	password, err := s.authRepo.Login(username)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (s *authService) CheckID(id int) bool {
	return s.authRepo.CheckID(id)
}

func (s *authService) Delete(id int) error {
	if err := s.authRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *authService) GetUserByEmail(email string) (*models.User, error) {
	return s.authRepo.GetUserByEmail(email)
}

func (s *authService) SaveRecoveryTokenToDB(user *models.User, token string) error {
	return s.authRepo.SaveRecoverToken(user, token)
}
