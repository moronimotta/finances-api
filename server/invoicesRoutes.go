package server

import (
	"finances-api/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initInvoicesRoutes() {

	s.app.POST("/invoices", func(ctx *gin.Context) {
		invoice := entities.Invoices{}

		if err := ctx.ShouldBindJSON(&invoice); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.pgHandler.Repository.CreateInvoice(&invoice); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating invoice"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice created successfully"})
	})
	s.app.GET("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		invoice, err := s.pgHandler.Repository.GetInvoiceByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}
		ctx.JSON(http.StatusOK, invoice)
	})
	s.app.PUT("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		invoice := entities.Invoices{}

		if err := ctx.ShouldBindJSON(&invoice); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		invoice.ID = id // Ensure the ID is set for the update

		if err := s.pgHandler.Repository.UpdateInvoice(&invoice); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating invoice"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully"})
	})
	s.app.DELETE("/invoices/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.pgHandler.Repository.DeleteInvoice(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting invoice"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
	})
	s.app.GET("/invoices/user/:user_id/product/:product_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")
		productID := ctx.Param("product_id")

		invoices, err := s.pgHandler.Repository.GetInvoicesByUserIDAndProductID(userID, productID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices", func(ctx *gin.Context) {
		invoices, err := s.pgHandler.Repository.GetAllInvoices()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}
		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/product/:product_id", func(ctx *gin.Context) {
		productID := ctx.Param("product_id")

		invoices, err := s.pgHandler.Repository.GetInvoicesByProductID(productID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/customer/:customer_id", func(ctx *gin.Context) {
		customerID := ctx.Param("customer_id")

		invoices, err := s.pgHandler.Repository.GetInvoicesByCustomerID(customerID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/payment_status/:payment_status", func(ctx *gin.Context) {
		paymentStatus := ctx.Param("payment_status")

		invoices, err := s.pgHandler.Repository.GetInvoicesByPaymentStatus(paymentStatus)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	})
	s.app.GET("/invoices/payment_method/:payment_method", func(ctx *gin.Context) {
		paymentMethod := ctx.Param("payment_method")

		invoices, err := s.pgHandler.Repository.GetInvoicesByPaymentMethod(paymentMethod)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting invoices"})
			return
		}

		ctx.JSON(http.StatusOK, invoices)
	})

}
