package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/delivery/http/v1/requests"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/delivery/http/v1/responses"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/services"
	"log/slog"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	Logger         *slog.Logger
	BookingService services.IBookingService
}

func (h *BookingHandler) GetHotels(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	sizeStr := ctx.Query("size")

	page, _ := strconv.Atoi(pageStr)
	if page == 0 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 10
	}

	hotels, err := h.BookingService.GetHotels(ctx.Request.Context(), page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{Err: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}

func (h *BookingHandler) GetHotelByUid(ctx *gin.Context) {
	uidStr := ctx.Param("hotelUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{Err: "Invalid UUID"})
		return
	}

	hotel, err := h.BookingService.GetHotelById(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, responses.ErrorResponse{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, hotel)
}

func (h *BookingHandler) CreateReservation(ctx *gin.Context) {
	var req requests.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{Err: err.Error()})
		return
	}

	res, err := h.BookingService.BookHotel(ctx.Request.Context(), req.HotelUid, req.Username, req.PaymentUid, req.HotelId, req.StartDate, req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (h *BookingHandler) GetUserReservations(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{Err: "Username required"})
		return
	}

	res, err := h.BookingService.GetReservations(ctx.Request.Context(), username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *BookingHandler) GetReservation(ctx *gin.Context) {
	uidStr := ctx.Param("reservationUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{Err: "Invalid UUID"})
		return
	}

	res, err := h.BookingService.GetReservationByUid(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, responses.ErrorResponse{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *BookingHandler) CancelReservation(ctx *gin.Context) {
	uidStr := ctx.Param("reservationUid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{Err: "Invalid UUID"})
		return
	}

	err = h.BookingService.CancelReservation(ctx.Request.Context(), uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{Err: err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
