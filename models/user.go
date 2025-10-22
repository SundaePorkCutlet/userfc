package models

type RegisterParameter struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);not null;unique;index" json:"email"`
	Password string `gorm:"type:text;not null" json:"password"`
	Role     string `gorm:"type:varchar(20);default:'user'" json:"role"`
}

type LoginParameter struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
