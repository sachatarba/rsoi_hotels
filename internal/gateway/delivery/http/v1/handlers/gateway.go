package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/services"
	"github.com/sachatarba/rsoi_hotels/pkg/circuitbreaker"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type GatewayHandler struct {
	service services.IGatewayService
	logger  *slog.Logger
}

func NewGatewayHandler(service services.IGatewayService, logger *slog.Logger) *GatewayHandler {
	return &GatewayHandler{
		service: service,
		logger:  logger,
	}
}

func isServiceUnavailable(err error) bool {
	if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
		return true
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "context deadline exceeded")
}

func (h *GatewayHandler) handleServiceError(c *gin.Context, err error) {
	if isServiceUnavailable(err) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Service is unavailable"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}

func (h *GatewayHandler) GetHotels(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size == 0 {
		size = 10
	}

	res, err := h.service.GetHotels(c.Request.Context(), page, size)
	if err != nil {
		h.logger.Error("GetHotels error", "error", err)
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) GetUserReservations(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}

	res, err := h.service.GetUserReservations(c.Request.Context(), username)
	if err != nil {
		h.logger.Error("GetUserReservations error", "error", err)
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) GetUserInfo(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}

	res, err := h.service.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		h.logger.Error("GetUserInfo error", "error", err)
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) GetReservation(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}
	uidStr := c.Param("reservationUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID"})
		return
	}

	res, err := h.service.GetReservation(c.Request.Context(), username, uid)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) BookHotel(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}

	var req entity.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	res, err := h.service.BookHotel(c.Request.Context(), username, req)
	if err != nil {
		h.logger.Error("BookHotel error", "error", err)
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GatewayHandler) CancelReservation(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}
	uidStr := c.Param("reservationUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID"})
		return
	}

	err = h.service.CancelReservation(c.Request.Context(), username, uid)
	if err != nil {
		h.logger.Error("CancelReservation error", "error", err)
		h.handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *GatewayHandler) GetLoyalty(c *gin.Context) {
	username := c.GetHeader("X-User-Name")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "X-User-Name header is required"})
		return
	}

	res, err := h.service.GetLoyalty(c.Request.Context(), username)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
