package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/v1/requests"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/v1/responses"
	"log/slog"
	"net/http"

	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/services"
)

type LoyaltyHandler struct {
	Service services.ILoyaltyService
	Logger  *slog.Logger
}

func (h *LoyaltyHandler) GetLoyalty(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Err: errors.New("username is required").Error(),
		})

		return
	}

	context := ctx.Request.Context()

	loyalty, err := h.Service.GetLoyaltyByUsername(context, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Err: err.Error(),
		})
		return
	}

	loyaltyResponse := responses.LoyaltyResponse{
		Id:               loyalty.Id,
		Username:         loyalty.Username,
		ReservationCount: loyalty.ReservationCount(),
		Status:           loyalty.Status().String(),
		Discount:         loyalty.Discount(),
	}

	ctx.JSON(http.StatusOK, loyaltyResponse)
}

func (h *LoyaltyHandler) AddReservations(ctx *gin.Context) {
	var reservationRequest requests.ReservationRequest
	err := ctx.ShouldBindJSON(&reservationRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Err: fmt.Errorf("request validation error: %w", err).Error(),
		})

		return
	}

	context := ctx.Request.Context()

	err = h.Service.AddReservationCount(context,
		reservationRequest.Username, reservationRequest.ReservationCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Err: fmt.Errorf("internal server error: %w", err).Error(),
		})
	}

	ctx.Status(http.StatusNoContent)
}
