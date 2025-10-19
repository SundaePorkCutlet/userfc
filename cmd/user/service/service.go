package service

import (
	"userfc/cmd/user/repository"
	"userfc/models"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(user *models.User) (int64, error) {
	userId, err := s.UserRepo.InsertNewUser(user)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
