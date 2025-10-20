package utils

import (
	"time"
	"userfc/config"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(config.GetJwtSecret()))
}
