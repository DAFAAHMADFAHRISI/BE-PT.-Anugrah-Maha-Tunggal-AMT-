package handler

import (
	"net/http"
	"strconv"

	"erp-anugrah-maha-tunggal/internal/domain"
	"erp-anugrah-maha-tunggal/internal/service"
	"erp-anugrah-maha-tunggal/pkg/response"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	masterService service.MasterService
}

func NewMasterHandler(masterService service.MasterService) *MasterHandler {
	return &MasterHandler{masterService: masterService}
}

// ================= CUSTOMERS =================
func (h *MasterHandler) GetCustomers(c *gin.Context) {
	data, err := h.masterService.GetAllCustomers()
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil data pelanggan", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil data pelanggan", data)
}

func (h *MasterHandler) CreateCustomer(c *gin.Context) {
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		response.BadRequest(c, "Input data pelanggan tidak valid", err.Error())
		return
	}

	if err := h.masterService.CreateCustomer(&customer); err != nil {
		response.InternalServerError(c, "Gagal menambahkan pelanggan", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Pelanggan berhasil ditambahkan", customer)
}

func (h *MasterHandler) UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "ID pelanggan tidak valid", nil)
		return
	}

	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		response.BadRequest(c, "Input data tidak valid", err.Error())
		return
	}
	customer.ID = uint(id)

	if err := h.masterService.UpdateCustomer(&customer); err != nil {
		response.InternalServerError(c, "Gagal memperbarui pelanggan", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Data pelanggan berhasil diperbarui", customer)
}

func (h *MasterHandler) DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	if err := h.masterService.DeleteCustomer(uint(id)); err != nil {
		response.InternalServerError(c, "Gagal menghapus pelanggan", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Pelanggan berhasil dihapus", nil)
}

// ================= FORKLIFT UNITS =================
func (h *MasterHandler) GetUnits(c *gin.Context) {
	statusFilter := c.Query("status")
	data, err := h.masterService.GetAllUnits(statusFilter)
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil data unit forklift", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil data unit forklift", data)
}

func (h *MasterHandler) CreateUnit(c *gin.Context) {
	var unit domain.ForkliftUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		response.BadRequest(c, "Input data unit forklift tidak valid", err.Error())
		return
	}

	if err := h.masterService.CreateUnit(&unit); err != nil {
		response.InternalServerError(c, "Gagal menambahkan unit forklift", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Unit forklift berhasil ditambahkan", unit)
}

func (h *MasterHandler) UpdateUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var unit domain.ForkliftUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		response.BadRequest(c, "Input data tidak valid", err.Error())
		return
	}
	unit.ID = uint(id)

	if err := h.masterService.UpdateUnit(&unit); err != nil {
		response.InternalServerError(c, "Gagal memperbarui unit", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Data unit forklift berhasil diperbarui", unit)
}

func (h *MasterHandler) DeleteUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	if err := h.masterService.DeleteUnit(uint(id)); err != nil {
		response.InternalServerError(c, "Gagal menghapus unit forklift", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Unit forklift berhasil dihapus", nil)
}

// ================= OPERATORS =================
func (h *MasterHandler) GetOperators(c *gin.Context) {
	data, err := h.masterService.GetAllOperators()
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil data operator", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil data operator", data)
}

func (h *MasterHandler) CreateOperator(c *gin.Context) {
	var operator domain.Operator
	if err := c.ShouldBindJSON(&operator); err != nil {
		response.BadRequest(c, "Input data operator tidak valid", err.Error())
		return
	}

	if err := h.masterService.CreateOperator(&operator); err != nil {
		response.InternalServerError(c, "Gagal menambahkan operator", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Operator berhasil ditambahkan", operator)
}
