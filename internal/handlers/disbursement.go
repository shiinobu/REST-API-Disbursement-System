package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"rest-api-disbursement-system/internal/middleware"
	"rest-api-disbursement-system/internal/services"
)

type DisbursementHandler struct {
	disbursements services.DisbursementService
}

type createDisbursementRequest struct {
	BeneficiaryName string  `json:"beneficiary_name" binding:"required,min=2,max=100"`
	BankName        string  `json:"bank_name" binding:"required,min=2,max=100"`
	AccountNumber   string  `json:"account_number" binding:"required,numeric,min=6,max=30"`
	Amount          float64 `json:"amount" binding:"gt=0"`
	Description     string  `json:"description" binding:"omitempty,max=255"`
}

type rejectDisbursementRequest struct {
	Reason string `json:"reason" binding:"required,min=5,max=255"`
}

func NewDisbursementHandler(disbursements services.DisbursementService) *DisbursementHandler {
	return &DisbursementHandler{disbursements: disbursements}
}

func (h *DisbursementHandler) MountRoutes(router *gin.RouterGroup) {
	router.GET("", h.List)
	router.POST("", h.Create)
	router.GET("/:id", h.Detail)
	router.PATCH("/:id/approve", h.Approve)
	router.PATCH("/:id/reject", h.Reject)
}

func (h *DisbursementHandler) List(c *gin.Context) {
	disbursements, err := h.disbursements.List()
	if err != nil {
		Error(c, http.StatusInternalServerError, "Gagal mengambil data disbursement", nil)
		return
	}

	Success(c, http.StatusOK, "Data disbursement berhasil diambil", disbursements)
}

func (h *DisbursementHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		Error(c, http.StatusUnauthorized, "Token tidak valid", nil)
		return
	}

	var request createDisbursementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ValidationError(c, err, request)
		return
	}

	disbursement, err := h.disbursements.Create(services.CreateDisbursementInput{
		RequesterID:     userID,
		BeneficiaryName: request.BeneficiaryName,
		BankName:        request.BankName,
		AccountNumber:   request.AccountNumber,
		Amount:          request.Amount,
		Description:     request.Description,
	})
	if err != nil {
		Error(c, http.StatusInternalServerError, "Gagal membuat disbursement", nil)
		return
	}

	Success(c, http.StatusCreated, "Disbursement berhasil dibuat", disbursement)
}

func (h *DisbursementHandler) Detail(c *gin.Context) {
	id, ok := disbursementID(c)
	if !ok {
		return
	}

	disbursement, err := h.disbursements.Detail(id)
	if err != nil {
		handleDisbursementError(c, err)
		return
	}

	Success(c, http.StatusOK, "Detail disbursement berhasil diambil", disbursement)
}

func (h *DisbursementHandler) Approve(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		Error(c, http.StatusUnauthorized, "Token tidak valid", nil)
		return
	}

	id, ok := disbursementID(c)
	if !ok {
		return
	}

	disbursement, err := h.disbursements.Approve(id, userID)
	if err != nil {
		handleDisbursementError(c, err)
		return
	}

	Success(c, http.StatusOK, "Disbursement berhasil disetujui", disbursement)
}

func (h *DisbursementHandler) Reject(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		Error(c, http.StatusUnauthorized, "Token tidak valid", nil)
		return
	}

	id, ok := disbursementID(c)
	if !ok {
		return
	}

	var request rejectDisbursementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ValidationError(c, err, request)
		return
	}

	disbursement, err := h.disbursements.Reject(id, userID, request.Reason)
	if err != nil {
		handleDisbursementError(c, err)
		return
	}

	Success(c, http.StatusOK, "Disbursement berhasil ditolak", disbursement)
}

func currentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

func disbursementID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		Error(c, http.StatusBadRequest, "ID disbursement tidak valid", gin.H{"id": "id harus berupa angka positif"})
		return 0, false
	}

	return uint(id), true
}

func handleDisbursementError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrDisbursementNotFound):
		Error(c, http.StatusNotFound, "Disbursement tidak ditemukan", gin.H{"id": "disbursement tidak ditemukan"})
	case errors.Is(err, services.ErrDisbursementAlreadyProcessed):
		Error(c, http.StatusConflict, "Disbursement sudah diproses", gin.H{"status": "hanya disbursement pending yang dapat diproses"})
	default:
		Error(c, http.StatusInternalServerError, "Gagal memproses disbursement", nil)
	}
}
