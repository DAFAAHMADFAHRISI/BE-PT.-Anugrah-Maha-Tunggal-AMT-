package handler

import (
	"net/http"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/service"
	"erp-anugrah-maha-tunggal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string          `json:"name" binding:"required"`
	Username string          `json:"username" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Role     domain.UserRole `json:"role" binding:"required"`
	Phone    string          `json:"phone"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Username dan password wajib diisi", err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"phone":    user.Phone,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data pendaftaran tidak valid", err.Error())
		return
	}

	user := domain.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Phone:    req.Phone,
		IsActive: true,
	}

	if err := h.authService.Register(&user, req.Password); err != nil {
		response.BadRequest(c, "Gagal mendaftarkan pengguna baru", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Pengguna berhasil didaftarkan", user)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Sesi tidak valid")
		return
	}

	user, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, "Profil pengguna tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Profil berhasil dimuat", user)
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil daftar pengguna", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar pengguna berhasil dimuat", users)
}

