package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initAdminRoutes() {

	s.app.GET("/admin/dashboard/info", func(ctx *gin.Context) {

		// get all products
		products, err := s.usecases.Db.GetAllProducts()
		if err != nil {
			slog.Error("Failed to get all products", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting products"})
			return
		}
		// get all transactions
		transactions, err := s.usecases.Db.GetAllTransactions()
		if err != nil {
			slog.Error("Error getting all transactions", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"products":     products,
			"transactions": transactions,
		})

	})

}
