package repository

import (
	"context"
	"userfc/models"

	"gorm.io/gorm"
)

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.Database.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.User{}, nil // 빈 사용자 객체 반환
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) InsertNewUser(ctx context.Context, user *models.User) (int64, error) {
	err := r.Database.WithContext(ctx).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepository) FindByUserId(ctx context.Context, userId int64) (*models.User, error) {
	var user models.User
	err := r.Database.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user, nil
		}
		return nil, err
	}
	return &user, nil
}
