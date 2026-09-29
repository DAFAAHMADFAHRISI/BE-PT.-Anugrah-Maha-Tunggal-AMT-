package handler

import (
	"net/http"

	"erp-anugrah-maha-tunggal/internal/service"
	"erp-anugrah-maha-tunggal/pkg/response"

	"github.com/gin-gonic/gin"
)

// DashboardHandler menangani semua endpoint dashboard berdasarkan role
type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetDirekturDashboard mengambil semua data ringkasan untuk dashboard eksekutif direktur
// GET /api/v1/dashboard/direktur
// Role yang diizinkan: direktur
func (h *DashboardHandler) GetDirekturDashboard(c *gin.Context) {
	data, err := h.dashboardService.GetDirekturDashboardData()
	if err != nil {
		response.InternalServerError(c, "Gagal memuat data dashboard direktur", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Data dashboard eksekutif berhasil dimuat", data)
}

// GetCurrentUserInfo mengambil info user yang sedang login berdasarkan JWT token
// GET /api/v1/auth/me
// Digunakan oleh frontend untuk routing otomatis ke dashboard sesuai role
func (h *DashboardHandler) GetCurrentUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	name, _ := c.Get("name")

	response.Success(c, http.StatusOK, "Informasi pengguna aktif berhasil dimuat", gin.H{
		"id":       userID,
		"username": username,
		"role":     role,
		"name":     name,
	})
}
