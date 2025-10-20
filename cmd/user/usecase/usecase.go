package usecase

import (
	"context"
	"errors"
	"userfc/cmd/user/service"
	"userfc/infrastructure/log"
	"userfc/models"
	"userfc/utils"
)

type UserUsecase struct {
	UserService service.UserService
}

func NewUserUsecase(userService service.UserService) *UserUsecase {
	return &UserUsecase{UserService: userService}
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUserByUserId(ctx context.Context, userId int64) (*models.User, error) {
	user, err := u.UserService.GetUserByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error hashing password: %s", err.Error())
		return err
	}
	user.Password = hashedPassword

	_, err = u.UserService.CreateUser(ctx, user)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error creating user: %s", err.Error())
		return err
	}
	return nil
}

func (u *UserUsecase) Login(ctx context.Context, loginParam *models.LoginParameter) (string, error) {
	user, err := u.UserService.GetUserByEmail(ctx, loginParam.Email)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error getting user by email: %s", err.Error())
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("user not found")
	}

	isMatch, err := utils.VerifyPassword(loginParam.Password, user.Password)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error verifying password: %s", err.Error())
		return "", err
	}
	if !isMatch {
		return "", errors.New("invalid password")
	}

	tokenString, err := utils.GenerateToken(user.ID)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error generating token: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
