package handler

import (
	"net/http"
	"strconv"

	"erp-anugrah-maha-tunggal/internal/service"
	"erp-anugrah-maha-tunggal/pkg/response"
	"github.com/gin-gonic/gin"
)

type SAWHandler struct {
	sawService service.SAWService
}

func NewSAWHandler(sawService service.SAWService) *SAWHandler {
	return &SAWHandler{sawService: sawService}
}

// GetLatestRanking GET /api/v1/saw/ranking
func (h *SAWHandler) GetLatestRanking(c *gin.Context) {
	result, err := h.sawService.GetLatestRanking(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil rekomendasi ranking SAW", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil mengambil rekomendasi prioritas pemeliharaan SAW", result)
}

// CalculateRanking POST /api/v1/saw/calculate
func (h *SAWHandler) CalculateRanking(c *gin.Context) {
	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	result, err := h.sawService.CalculateAndSave(c.Request.Context(), userID, body.Notes)
	if err != nil {
		response.BadRequest(c, "Gagal mengeksekusi kalkulasi SAW", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Kalkulasi algoritma SAW berhasil dieksekusi dan disimpan", result)
}

// GetCriteriaAndWeights GET /api/v1/saw/criteria
func (h *SAWHandler) GetCriteriaAndWeights(c *gin.Context) {
	kriteria, bobot, err := h.sawService.GetCriteriaAndWeights(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil kriteria dan bobot", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data kriteria dan bobot preferensi SAW", gin.H{
		"kriteria": kriteria,
		"bobot":    bobot,
	})
}

// UpdateWeights PUT /api/v1/saw/weights
func (h *SAWHandler) UpdateWeights(c *gin.Context) {
	var req struct {
		Weights []service.WeightInputRequest `json:"weights" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format data bobot tidak valid", err.Error())
		return
	}

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	if err := h.sawService.UpdateWeights(c.Request.Context(), req.Weights, userID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Trigger hitung ulang otomatis setelah bobot diperbarui
	calcResult, _ := h.sawService.CalculateAndSave(c.Request.Context(), userID, "Kalkulasi otomatis pasca-pembaruan bobot preferensi")

	response.Success(c, http.StatusOK, "Bobot preferensi kriteria berhasil diperbarui dan peringkat diselaraskan", calcResult)
}

// GetRankingHistory GET /api/v1/saw/history
func (h *SAWHandler) GetRankingHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	history, err := h.sawService.GetRankingHistory(c.Request.Context(), limit)
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil riwayat ranking SAW", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil riwayat snapshot ranking SAW", history)
}

// CreateOperationalLog POST /api/v1/operational-logs
func (h *SAWHandler) CreateOperationalLog(c *gin.Context) {
	var req service.OperationalLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data input log tidak valid: "+err.Error(), nil)
		return
	}

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	entry, err := h.sawService.CreateOperationalLog(c.Request.Context(), req, userID)
	if err != nil {
		response.InternalServerError(c, "Gagal menyimpan log operasional", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Log operasional berhasil disimpan dan peringkat SAW diperbarui", entry)
}

// GetOperationalLogs GET /api/v1/operational-logs
func (h *SAWHandler) GetOperationalLogs(c *gin.Context) {
	var unitIDPtr *uint
	if unitIDStr := c.Query("unit_id"); unitIDStr != "" {
		if u, err := strconv.ParseUint(unitIDStr, 10, 32); err == nil {
			uid := uint(u)
			unitIDPtr = &uid
		}
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.sawService.GetOperationalLogs(c.Request.Context(), unitIDPtr, limit)
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil log operasional", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar log operasional armada", logs)
}

// CreateSPKMaintenance POST /api/v1/maintenance/spk
func (h *SAWHandler) CreateSPKMaintenance(c *gin.Context) {
	var req service.SPKMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Data form SPK tidak valid: "+err.Error(), nil)
		return
	}

	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	spk, err := h.sawService.CreateSPKMaintenance(c.Request.Context(), req, userID)
	if err != nil {
		response.InternalServerError(c, "Gagal menerbitkan SPK", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Surat Perintah Kerja (SPK) Pemeliharaan berhasil diterbitkan dan status unit diubah ke MAINTENANCE", spk)
}

// GetSPKMaintenanceList GET /api/v1/maintenance/spk
func (h *SAWHandler) GetSPKMaintenanceList(c *gin.Context) {
	list, err := h.sawService.GetSPKMaintenanceList(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Gagal mengambil daftar SPK", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar Surat Perintah Kerja (SPK) Pemeliharaan", list)
}

// GetSPKMaintenanceByID GET /api/v1/maintenance/spk/:id
func (h *SAWHandler) GetSPKMaintenanceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "ID SPK tidak valid", nil)
		return
	}

	spk, err := h.sawService.GetSPKMaintenanceByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Dokumen SPK tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil rincian dokumen SPK Pemeliharaan", spk)
}
