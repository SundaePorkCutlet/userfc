package handler

import (
	"net/http"
	"userfc/cmd/user/usecase"
	"userfc/infrastructure/log"

	"userfc/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUsecase}
}

func (h *UserHandler) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var registerParam models.RegisterParameter
	if err := c.ShouldBindJSON(&registerParam); err != nil {
		log.Logger.Info().Err(err).Msg("Invalid JSON format in registration request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(registerParam.Password) < 8 {
		log.Logger.Info().Msg("Password too short - less than 8 characters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	if registerParam.Password != registerParam.ConfirmPassword {
		log.Logger.Info().Msg("Password and ConfirmPassword do not match")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password and ConfirmPassword do not match"})
		return
	}

	user, err := h.UserUsecase.GetUserByEmail(registerParam.Email)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error getting user by email")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	err = h.UserUsecase.RegisterUser(&models.User{
		Name:     registerParam.Name,
		Email:    registerParam.Email,
		Password: registerParam.Password, //plain text password
	})
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error registering user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
