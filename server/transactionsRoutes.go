package server

import (
	"finances-api/entities"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initTransactionsRoutes() {

	s.app.POST("/transactions", func(ctx *gin.Context) {
		transaction := entities.Transactions{}

		if err := ctx.ShouldBindJSON(&transaction); err != nil {
			slog.Error("Invalid request for creating transaction", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.usecases.Db.CreateTransaction(&transaction); err != nil {
			slog.Error("Error creating transaction", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating transaction"})
			return
		}

		slog.Info("Transaction created successfully", "transaction_id", transaction.ID)
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
	})

	s.app.GET("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		transaction, err := s.usecases.Db.GetTransactionByID(id)
		if err != nil {
			slog.Error("Transaction not found", "id", id, "error", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		slog.Info("Transaction retrieved successfully", "id", id)
		ctx.JSON(http.StatusOK, transaction)
	})

	s.app.PUT("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		transaction := entities.Transactions{}

		if err := ctx.ShouldBindJSON(&transaction); err != nil {
			slog.Error("Invalid request for updating transaction", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		transaction.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateTransaction(&transaction); err != nil {
			slog.Error("Error updating transaction", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating transaction"})
			return
		}

		slog.Info("Transaction updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
	})
	s.app.DELETE("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.usecases.Db.DeleteTransaction(id); err != nil {
			slog.Error("Error deleting transaction", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting transaction"})
			return
		}
		slog.Info("Transaction deleted successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
	})
	s.app.GET("/transactions/user/:user_id/product/:product_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")
		productID := ctx.Param("product_id")

		transactions, err := s.usecases.Db.GetTransactionsByUserIDAndProductID(userID, productID)
		if err != nil {
			slog.Error("Error getting transactions by user and product", "user_id", userID, "product_id", productID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		slog.Info("Transactions retrieved successfully", "user_id", userID, "product_id", productID, "count", len(transactions))
		ctx.JSON(http.StatusOK, transactions)
	})

	s.app.GET("/transactions", func(ctx *gin.Context) {
		transactions, err := s.usecases.Db.GetAllTransactions()
		if err != nil {
			slog.Error("Error getting all transactions", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		slog.Info("All transactions retrieved successfully", "count", len(transactions))
		ctx.JSON(http.StatusOK, transactions)
	})

	s.app.GET("/transactions/product/:product_id", func(ctx *gin.Context) {
		productID := ctx.Param("product_id")

		transactions, err := s.usecases.Db.GetTransactionsByProductID(productID)
		if err != nil {
			slog.Error("Error getting transactions by product", "product_id", productID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		slog.Info("Transactions retrieved successfully", "product_id", productID, "count", len(transactions))
		ctx.JSON(http.StatusOK, transactions)
	})

}
