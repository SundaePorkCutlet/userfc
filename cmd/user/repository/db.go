package repository

import "userfc/models"

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.Database.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) InsertNewUser(user *models.User) (int64, error) {
	err := r.Database.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
