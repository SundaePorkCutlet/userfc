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

	user, err := h.UserUsecase.GetUserByEmail(c.Request.Context(), registerParam.Email)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error getting user by email")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	err = h.UserUsecase.RegisterUser(c.Request.Context(), &models.User{
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

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIdString, ok := c.Get("user_id")
	if !ok {
		log.Logger.Info().Msg("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userId, ok := userIdString.(float64)
	if !ok {
		log.Logger.Info().Msg("Invalid user ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserUsecase.GetUserByUserId(c.Request.Context(), int64(userId))
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error getting user by user id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.ID == 0 {
		log.Logger.Info().Msg("User not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": user.Name, "email": user.Email})

}

func (h *UserHandler) Login(c *gin.Context) {
	var loginParam models.LoginParameter
	if err := c.ShouldBindJSON(&loginParam); err != nil {
		log.Logger.Info().Err(err).Msgf("Invalid JSON format in login request: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(loginParam.Password) < 8 {
		log.Logger.Info().Msg("Password too short - less than 8 characters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	token, err := h.UserUsecase.Login(c.Request.Context(), &loginParam)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error logging in: %s", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
