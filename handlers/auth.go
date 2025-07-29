package handlers

import (
	"database/sql"
	"go-orm-vs-sql/models"
	"go-orm-vs-sql/repositories"
	"go-orm-vs-sql/types"
	"go-orm-vs-sql/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo repositories.UserRepository
}

func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// ユーザー登録
func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// ユーザー重複チェック
	emailExists, err := h.userRepo.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Database error"})
		return
	}
	
	usernameExists, err := h.userRepo.UsernameExists(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Database error"})
		return
	}
	
	if emailExists || usernameExists {
		c.JSON(http.StatusConflict, types.ErrorResponse{Error: "User already exists"})
		return
	}

	// パスワードハッシュ化
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to hash password"})
		return
	}

	// ユーザー作成
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create user"})
		return
	}

	// JWT生成
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, types.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// ログイン
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// ユーザー検索
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Database error"})
		return
	}

	// パスワード検証
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// JWT生成
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, types.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// 現在のユーザー情報取得
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not found"})
		return
	}

	user, err := h.userRepo.GetByID(userID.(uint))
	if err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}