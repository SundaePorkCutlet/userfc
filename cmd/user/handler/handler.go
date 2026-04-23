package handler

import (
	"context"
	"net/http"
	"strings"
	"time"
	"userfc/cmd/user/usecase"
	"userfc/config"
	"userfc/infrastructure/log"
	"userfc/infrastructure/tokenblacklist"

	"userfc/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
	Blacklist   *tokenblacklist.TokenBlacklist
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUsecase}
}

func (h *UserHandler) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}

// Register godoc
// @Summary 회원가입
// @Description 새로운 사용자를 등록합니다.
// @Tags USER
// @Accept json
// @Produce json
// @Param body body models.RegisterParameter true "회원가입 요청"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/register [post]
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

// GetUserInfo godoc
// @Summary 내 정보 조회
// @Description JWT 인증된 사용자의 이름/이메일을 조회합니다.
// @Tags USER
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user-info [get]
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

// Login godoc
// @Summary 로그인
// @Description 이메일/비밀번호로 로그인 후 JWT 토큰을 발급합니다.
// @Tags USER
// @Accept json
// @Produce json
// @Param body body models.LoginParameter true "로그인 요청"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/login [post]
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

// Logout godoc
// @Summary 로그아웃
// @Description 현재 JWT 토큰을 블랙리스트 처리하여 즉시 무효화합니다.
// @Tags USER
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorization header"})
		return
	}

	rawToken := parts[1]

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetJwtSecret()), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token claims"})
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing expiration in token"})
		return
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if ttl <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "token already expired"})
		return
	}

	if err := h.Blacklist.Add(context.Background(), rawToken, ttl); err != nil {
		log.Logger.Info().Err(err).Msg("Failed to blacklist token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
