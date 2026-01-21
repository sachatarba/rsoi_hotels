package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/services"
	"log/slog"
	"net/http"
)

type PaymentHandler struct {
	Service services.IPaymentService
	Logger  *slog.Logger
}

type CreatePaymentRequest struct {
	Price int `json:"price"`
}

func (h *PaymentHandler) Create(ctx *gin.Context) {
	var req CreatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, err := h.Service.Create(ctx.Request.Context(), req.Price)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"payment_uid": uid})
}

func (h *PaymentHandler) Cancel(ctx *gin.Context) {
	uidStr := ctx.Param("paymentUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	err = h.Service.Cancel(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *PaymentHandler) Get(ctx *gin.Context) {
	uidStr := ctx.Param("paymentUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	payment, err := h.Service.GetDetails(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
