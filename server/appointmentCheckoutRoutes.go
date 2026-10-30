package server

import (
	"errors"
	"finances-api/entities"
	"finances-api/usecases"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initAppointmentCheckoutRoutes() {
	s.app.POST("/checkout/appointments", barberServiceToken(), func(ctx *gin.Context) {
		var input entities.AppointmentCheckoutRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment_id and at least one service_id are required"})
			return
		}

		stripeUsecase := usecases.NewStripeUsecase(s.usecases)
		session, total, currency, err := stripeUsecase.CreateAppointmentCheckoutSession(input.AppointmentID, input.ServiceIDs)
		if err != nil {
			switch {
			case errors.Is(err, usecases.ErrInvalidAppointmentCheckout):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, usecases.ErrAppointmentServiceUnavailable):
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.Is(err, usecases.ErrAppointmentCheckoutExists):
				ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				slog.Error("Failed to create appointment checkout session", "error", err)
				ctx.JSON(http.StatusBadGateway, gin.H{"error": "Unable to start appointment checkout"})
			}
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"checkout_url":   session.URL,
			"session_id":     session.ID,
			"amount_total":   total,
			"currency":       currency,
			"payment_status": "pending",
		})
	})
}
