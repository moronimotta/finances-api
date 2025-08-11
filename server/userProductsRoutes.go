package server

import (
	"finances-api/entities"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initUserProductsRoutes() {

	s.app.POST("/user-products", func(ctx *gin.Context) {
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			slog.Error("Invalid request when creating user product", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.usecases.Db.CreateUserProduct(userProduct.UserID, userProduct.ProductID); err != nil {
			slog.Error("Error creating user product", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user product"})
			return
		}

		slog.Info("User product created successfully", "user_id", userProduct.UserID, "product_id", userProduct.ProductID)
		ctx.JSON(http.StatusOK, gin.H{"message": "User product created successfully"})
	})
	s.app.GET("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userProduct, err := s.usecases.Db.GetUserProductByID(id)
		if err != nil {
			slog.Error("User product not found", "id", id, "error", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User product not found"})
			return
		}
		slog.Info("User product retrieved successfully", "id", id)
		ctx.JSON(http.StatusOK, userProduct)
	})
	s.app.PUT("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			slog.Error("Invalid request when updating user product", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userProduct.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateUserProduct(&userProduct); err != nil {
			slog.Error("Error updating user product", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user product"})
			return
		}

		slog.Info("User product updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "User product updated successfully"})
	})

	s.app.PUT("/user-products/status/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			slog.Error("Invalid request when updating user product", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userProduct.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateUserProductStatus(&userProduct); err != nil {
			slog.Error("Error updating user product status", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user product status"})
			return
		}

		slog.Info("User product status updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "User product status updated successfully"})
	})
	s.app.DELETE("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.usecases.Db.DeleteUserProduct(id); err != nil {
			slog.Error("Error deleting user product", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user product"})
			return
		}

		slog.Info("User product deleted successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "User product deleted successfully"})
	})
	s.app.GET("/user-products/user/:user_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")

		userProducts, err := s.usecases.Db.GetUserProductsByUserID(userID)
		if err != nil {
			slog.Error("Error getting user products", "user_id", userID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user products"})
			return
		}
		slog.Info("User products retrieved successfully", "user_id", userID, "count", len(userProducts))
		ctx.JSON(http.StatusOK, userProducts)
	})

}
