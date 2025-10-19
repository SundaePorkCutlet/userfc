package usecase

import (
	"userfc/cmd/user/service"
	"userfc/models"
)

type UserUsecase struct {
	UserService service.UserService
}

func NewUserUsecase(userService service.UserService) *UserUsecase {
	return &UserUsecase{UserService: userService}
}

func (u *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.UserService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) RegisterUser(user *models.User) error {
	_, err := u.UserService.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
