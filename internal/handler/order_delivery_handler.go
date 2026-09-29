package handler

import (
	"net/http"
	"strconv"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/service"
	"erp-anugrah-maha-tunggal/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderDeliveryHandler struct {
	orderService    service.OrderService
	deliveryService service.DeliveryService
}

func NewOrderDeliveryHandler(
	orderService service.OrderService,
	deliveryService service.DeliveryService,
) *OrderDeliveryHandler {
	return &OrderDeliveryHandler{
		orderService:    orderService,
		deliveryService: deliveryService,
	}
}

// ================= RENTAL ORDERS =================
func (h *OrderDeliveryHandler) GetOrders(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.orderService.GetAllOrders(status)
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil daftar pesanan sewa", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar pesanan sewa berhasil dimuat", orders)
}

func (h *OrderDeliveryHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "ID pesanan tidak valid", nil)
		return
	}

	order, err := h.orderService.GetOrderByID(uint(id))
	if err != nil {
		response.NotFound(c, "Pesanan sewa tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Detail pesanan sewa berhasil dimuat", order)
}

func (h *OrderDeliveryHandler) CreateOrder(c *gin.Context) {
	var order domain.RentalOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		response.BadRequest(c, "Data pesanan tidak valid", err.Error())
		return
	}

	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		order.CreatedBy = &uid
	}

	if err := h.orderService.CreateOrder(&order); err != nil {
		response.InternalServerError(c, "Gagal membuat pesanan sewa", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Pesanan sewa berhasil dibuat", order)
}

func (h *OrderDeliveryHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Status baru wajib diisi", err.Error())
		return
	}

	if err := h.orderService.UpdateOrderStatus(uint(id), req.Status); err != nil {
		response.InternalServerError(c, "Gagal memperbarui status pesanan", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Status pesanan berhasil diperbarui", nil)
}

// ================= DELIVERY LETTERS (SURAT JALAN) =================
func (h *OrderDeliveryHandler) GetDeliveryLetters(c *gin.Context) {
	status := c.Query("status")
	letters, err := h.deliveryService.GetAllDeliveryLetters(status)
	if err != nil {
		response.InternalServerError(c, "Gagal memuat surat jalan", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Daftar surat jalan berhasil dimuat", letters)
}

func (h *OrderDeliveryHandler) GetDeliveryLetterByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "ID surat jalan tidak valid", nil)
		return
	}

	letter, err := h.deliveryService.GetDeliveryLetterByID(uint(id))
	if err != nil {
		response.NotFound(c, "Surat jalan tidak ditemukan")
		return
	}
	response.Success(c, http.StatusOK, "Detail surat jalan berhasil dimuat", letter)
}

func (h *OrderDeliveryHandler) CreateDeliveryLetter(c *gin.Context) {
	var letter domain.DeliveryLetter
	if err := c.ShouldBindJSON(&letter); err != nil {
		response.BadRequest(c, "Data surat jalan tidak valid", err.Error())
		return
	}

	if err := h.deliveryService.CreateDeliveryLetter(&letter); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Surat jalan berhasil diterbitkan", letter)
}

func (h *OrderDeliveryHandler) UpdateOperationalStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Status string `json:"status" binding:"required"` // ASSIGNED, ON_THE_WAY, WORKING, FINISHED
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Status operasional wajib diisi", err.Error())
		return
	}

	if err := h.deliveryService.UpdateOperationalStatus(uint(id), req.Status, req.Notes); err != nil {
		response.InternalServerError(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Status operasional berhasil diperbarui", nil)
}

func (h *OrderDeliveryHandler) GetUnitMovementLogs(c *gin.Context) {
	unitIDStr := c.Param("unit_id")
	unitID, _ := strconv.ParseUint(unitIDStr, 10, 32)

	logs, err := h.deliveryService.GetMovementLogs(uint(unitID))
	if err != nil {
		response.InternalServerError(c, "Gagal memuat log audit pergerakan unit", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Log pergerakan unit berhasil dimuat", logs)
}

