package service

import (
	"context"
	"userfc/cmd/user/repository"
	"userfc/models"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	userId, err := s.UserRepo.InsertNewUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *UserService) GetUserByUserId(ctx context.Context, userId int64) (*models.User, error) {
	user, err := s.UserRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
