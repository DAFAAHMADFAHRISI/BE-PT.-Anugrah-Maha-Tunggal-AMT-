package middleware

import (
	"strings"

	"erp-anugrah-maha-tunggal/pkg/response"
	"erp-anugrah-maha-tunggal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Akses ditolak: Token autentikasi tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Format header Authorization harus: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token tidak valid atau telah kedaluwarsa")
			c.Abort()
			return
		}

		// Simpan claims di context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("name", claims.Name)

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "Akses ditolak: Hak akses tidak teridentifikasi")
			c.Abort()
			return
		}

		userRole := roleVal.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if r == userRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			response.Unauthorized(c, "Akses ditolak: Anda tidak memiliki wewenang untuk modul ini")
			c.Abort()
			return
		}

		c.Next()
	}
}
