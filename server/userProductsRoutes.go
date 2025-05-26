package server

import (
	"finances-api/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initUserProductsRoutes() {

	s.app.POST("/user-products", func(ctx *gin.Context) {
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.pgHandler.Repository.CreateUserProduct(userProduct.UserID, userProduct.ProductID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user product"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User product created successfully"})
	})
	s.app.GET("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		userProduct, err := s.pgHandler.Repository.GetUserProductByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User product not found"})
			return
		}
		ctx.JSON(http.StatusOK, userProduct)
	})
	s.app.PUT("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userProduct.ID = id // Ensure the ID is set for the update

		if err := s.pgHandler.Repository.UpdateUserProduct(&userProduct); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user product"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User product updated successfully"})
	})

	s.app.PUT("/user-products/status/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var userProduct entities.UserProducts

		if err := ctx.ShouldBindJSON(&userProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userProduct.ID = id // Ensure the ID is set for the update

		if err := s.pgHandler.Repository.UpdateUserProductStatus(&userProduct); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user product status"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "User product status updated successfully"})
	})
	s.app.DELETE("/user-products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.pgHandler.Repository.DeleteUserProduct(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user product"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "User product deleted successfully"})
	})
	s.app.GET("/user-products/user/:user_id", func(ctx *gin.Context) {
		userID := ctx.Param("user_id")

		userProducts, err := s.pgHandler.Repository.GetUserProductsByUserID(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user products"})
			return
		}
		ctx.JSON(http.StatusOK, userProducts)
	})

}
