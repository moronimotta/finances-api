package server

import (
	"finances-api/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initTransactionsRoutes() {

	s.app.POST("/transactions", func(ctx *gin.Context) {
		transaction := entities.Transactions{}

		if err := ctx.ShouldBindJSON(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.usecases.Db.CreateTransaction(&transaction); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating transaction"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
	})

	s.app.GET("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		transaction, err := s.usecases.Db.GetTransactionByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		ctx.JSON(http.StatusOK, transaction)
	})

	s.app.PUT("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		transaction := entities.Transactions{}

		if err := ctx.ShouldBindJSON(&transaction); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		transaction.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateTransaction(&transaction); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating transaction"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
	})
	s.app.DELETE("/transactions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.usecases.Db.DeleteTransaction(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting transaction"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
	})
	s.app.GET("/transactions/user/:user_id/product/:product_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")
		productID := ctx.Param("product_id")

		transactions, err := s.usecases.Db.GetTransactionsByUserIDAndProductID(userID, productID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		ctx.JSON(http.StatusOK, transactions)
	})

	s.app.GET("/transactions", func(ctx *gin.Context) {
		transactions, err := s.usecases.Db.GetAllTransactions()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		ctx.JSON(http.StatusOK, transactions)
	})

	s.app.GET("/transactions/product/:product_id", func(ctx *gin.Context) {
		productID := ctx.Param("product_id")

		transactions, err := s.usecases.Db.GetTransactionsByProductID(productID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting transactions"})
			return
		}
		ctx.JSON(http.StatusOK, transactions)
	})

}
