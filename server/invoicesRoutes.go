package server

import (
	"finances-api/entities"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initInvoicesRoutes() {

	s.app.POST("/invoices", func(ctx *gin.Context) {
		invoice := entities.Invoices{}

		if err := ctx.ShouldBindJSON(&invoice); err != nil {
			slog.Error("Invalid request when creating invoice", "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.usecases.Db.CreateInvoice(&invoice); err != nil {
			slog.Error("Error creating invoice", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating invoice"})
			return
		}

		slog.Info("Invoice created successfully", "invoice_id", invoice.ID)
		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice created successfully"})
	})
	s.app.GET("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		invoice, err := s.usecases.Db.GetInvoiceByID(id)
		if err != nil {
			slog.Error("Invoice not found", "id", id, "error", err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		slog.Info("Invoice retrieved successfully", "id", id)
		ctx.JSON(http.StatusOK, invoice)
	})
	s.app.PUT("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		invoice := entities.Invoices{}

		if err := ctx.ShouldBindJSON(&invoice); err != nil {
			slog.Error("Invalid request when updating invoice", "id", id, "error", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		invoice.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateInvoice(&invoice); err != nil {
			slog.Error("Error updating invoice", "id", id, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating invoice"})
			return
		}

		slog.Info("Invoice updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully"})
	})
	s.app.DELETE("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.usecases.Db.DeleteInvoice(id); err != nil {
			slog.Error("Error deleting invoice", "id", id, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting invoice"})
			return
		}
		slog.Info("Invoice deleted successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
	})
	s.app.GET("/invoices/user/:user_id/product/:product_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")
		productID := ctx.Param("product_id")

		invoices, err := s.usecases.Db.GetInvoicesByUserIDAndProductID(userID, productID)
		if err != nil {
			slog.Error("Error getting invoices by user and product", "user_id", userID, "product_id", productID, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		slog.Info("Invoices retrieved successfully by user and product", "user_id", userID, "product_id", productID, "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices", func(ctx *gin.Context) {
		invoices, err := s.usecases.Db.GetAllInvoices()
		if err != nil {
			slog.Error("Error getting all invoices", "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}
		slog.Info("All invoices retrieved successfully", "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/product/:product_id", func(ctx *gin.Context) {
		productID := ctx.Param("product_id")

		invoices, err := s.usecases.Db.GetInvoicesByProductID(productID)
		if err != nil {
			slog.Error("Error getting invoices by product", "product_id", productID, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		slog.Info("Invoices retrieved successfully by product", "product_id", productID, "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/customer/:customer_id", func(ctx *gin.Context) {
		customerID := ctx.Param("customer_id")

		invoices, err := s.usecases.Db.GetInvoicesByCustomerID(customerID)
		if err != nil {
			slog.Error("Error getting invoices by customer", "customer_id", customerID, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		slog.Info("Invoices retrieved successfully by customer", "customer_id", customerID, "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/payment_status/:payment_status", func(ctx *gin.Context) {
		paymentStatus := ctx.Param("payment_status")

		invoices, err := s.usecases.Db.GetInvoicesByPaymentStatus(paymentStatus)
		if err != nil {
			slog.Error("Error getting invoices by payment status", "payment_status", paymentStatus, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		slog.Info("Invoices retrieved successfully by payment status", "payment_status", paymentStatus, "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/payment_method/:payment_method", func(ctx *gin.Context) {
		paymentMethod := ctx.Param("payment_method")

		invoices, err := s.usecases.Db.GetInvoicesByPaymentMethod(paymentMethod)
		if err != nil {
			slog.Error("Error getting invoices by payment method", "payment_method", paymentMethod, "error", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		slog.Info("Invoices retrieved successfully by payment method", "payment_method", paymentMethod, "count", len(invoices))
		ctx.JSON(http.StatusOK, invoices)
	})

}
