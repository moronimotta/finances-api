package server

import (
	"finances-api/entities"
	"finances-api/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initProductsRoutes() {
	s.app.POST("/products", func(ctx *gin.Context) {
		product := entities.Products{}

		if err := ctx.ShouldBindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		gatewayHandler, err := handlers.NewGatewayHttpHandler(product.GatewayName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}

		product.ExternalID, err = gatewayHandler.CreateProduct(product.Name, product.Description)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in gateway"})
			return
		}

		_, err = gatewayHandler.CreatePrice(product.ExternalID, product.Price, product.Currency)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating price in gateway"})
			return
		}

		s.pgHandler.Repository.CreateProduct(product.Name, product.Description, product.ExternalID, product.GatewayName, product.Price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in database"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	})

	s.app.GET("/products", func(ctx *gin.Context) {
		products, err := s.pgHandler.Repository.GetAllProducts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting products"})
			return
		}
		ctx.JSON(http.StatusOK, products)
	})
	s.app.GET("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		product, err := s.pgHandler.Repository.GetProductByID(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		ctx.JSON(http.StatusOK, product)
	})
	s.app.PUT("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		product := entities.Products{}
		if err := ctx.ShouldBindJSON(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		product.ID = id
		if err := s.pgHandler.Repository.UpdateProduct(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating product"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	})

	s.app.DELETE("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		product, err := s.pgHandler.Repository.GetProductByID(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		if err := s.pgHandler.Repository.DeleteProduct(product.ID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting product"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})

	s.app.GET("/products/g/:external_id", func(ctx *gin.Context) {
		externalID := ctx.Param("external_id")
		product, err := s.pgHandler.Repository.GetProductByExternalID(externalID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		ctx.JSON(http.StatusOK, product)
	})
}
