package main

import (
	"log"
	"os"

	"erp-anugrah-maha-tunggal/config"
	"erp-anugrah-maha-tunggal/internal/handler"
	"erp-anugrah-maha-tunggal/internal/middleware"
	"erp-anugrah-maha-tunggal/internal/repository"
	"erp-anugrah-maha-tunggal/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Database MySQL & GORM AutoMigrate
	db := config.InitDB()

	// 2. Inisialisasi Repositories
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	operatorRepo := repository.NewOperatorRepository(db)
	orderRepo := repository.NewRentalOrderRepository(db)
	deliveryRepo := repository.NewDeliveryLetterRepository(db)

	// 3. Inisialisasi Services
	authService := service.NewAuthService(userRepo)
	masterService := service.NewMasterService(customerRepo, unitRepo, operatorRepo)
	orderService := service.NewOrderService(orderRepo, unitRepo)
	deliveryService := service.NewDeliveryService(deliveryRepo, orderRepo, unitRepo, operatorRepo)
	dashboardService := service.NewDashboardService(unitRepo, orderRepo, customerRepo, deliveryRepo)

	// 4. Inisialisasi Handlers
	authHandler := handler.NewAuthHandler(authService)
	masterHandler := handler.NewMasterHandler(masterService)
	orderDeliveryHandler := handler.NewOrderDeliveryHandler(orderService, deliveryService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// 5. Setup Gin Router
	r := gin.Default()

	// Setup CORS agar frontend TypeScript/React dapat mengakses API
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"app":     "ERP PT. Anugrah Maha Tunggal",
			"version": "1.0.0",
		})
	})

	// Grouping API v1
	v1 := r.Group("/api/v1")
	{
		// Public Auth
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/register", authHandler.Register)
		}

		// Protected Endpoints (Perlu JWT Token)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Info pengguna aktif dari token — digunakan frontend untuk routing otomatis
			protected.GET("/auth/me", dashboardHandler.GetCurrentUserInfo)

			// Profil & Manajemen Akun Pengguna
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.GET("/users", authHandler.GetUsers)

			// Dashboard Direktur (Hanya role 'direktur')
			protected.GET("/dashboard/direktur",
				middleware.RequireRole("direktur"),
				dashboardHandler.GetDirekturDashboard,
			)

			// Master Pelanggan (Customers)
			protected.GET("/customers", masterHandler.GetCustomers)
			protected.POST("/customers", masterHandler.CreateCustomer)
			protected.PUT("/customers/:id", masterHandler.UpdateCustomer)
			protected.DELETE("/customers/:id", masterHandler.DeleteCustomer)

			// Master Unit Forklift
			protected.GET("/units", masterHandler.GetUnits)
			protected.POST("/units", masterHandler.CreateUnit)
			protected.PUT("/units/:id", masterHandler.UpdateUnit)
			protected.DELETE("/units/:id", masterHandler.DeleteUnit)

			// Master Operator
			protected.GET("/operators", masterHandler.GetOperators)
			protected.POST("/operators", masterHandler.CreateOperator)

			// Transaksi Pesanan Sewa (Rental Orders)
			protected.GET("/orders", orderDeliveryHandler.GetOrders)
			protected.GET("/orders/:id", orderDeliveryHandler.GetOrderByID)
			protected.POST("/orders", orderDeliveryHandler.CreateOrder)
			protected.PATCH("/orders/:id/status", orderDeliveryHandler.UpdateOrderStatus)

			// Surat Jalan Operasional (Delivery Letters)
			protected.GET("/delivery-letters", orderDeliveryHandler.GetDeliveryLetters)
			protected.GET("/delivery-letters/:id", orderDeliveryHandler.GetDeliveryLetterByID)
			protected.POST("/delivery-letters", orderDeliveryHandler.CreateDeliveryLetter)
			protected.PATCH("/delivery-letters/:id/operational-status", orderDeliveryHandler.UpdateOperationalStatus)

			// Log Pergerakan & Status Unit (Audit Trail)
			protected.GET("/units/:unit_id/logs", orderDeliveryHandler.GetUnitMovementLogs)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Aplikasi Backend ERP PT. Anugrah Maha Tunggal berjalan di port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
