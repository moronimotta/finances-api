package server

import (
	"finances-api/entities"
	"finances-api/handlers"
	"finances-api/utils/meta"
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

		product.GatewayProductExternalID, err = gatewayHandler.CreateProduct(product.Name, product.Description, product)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in gateway"})
			return
		}

		product.GatewayPriceExternalID, err = gatewayHandler.CreatePrice(product.GatewayProductExternalID, product.Price, product.Currency)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating price in gateway"})
			return
		}

		s.pgHandler.Repository.CreateProduct(&product)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in database"})
			return
		}

		meta := meta.New()
		meta.Add("local_product_id", product.ID)
		meta.Add("gateway_current_price_id", product.GatewayPriceExternalID)

		gatewayHandler.UpdateProduct(product.GatewayProductExternalID, "", "", meta)

		ctx.JSON(http.StatusOK, gin.H{"message": "Product created successfully"})
	})

	s.app.GET("/products", func(ctx *gin.Context) {
		products, err := s.pgHandler.Repository.GetAllProducts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting products"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"products": products})
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

		if product.Price != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be changed"})
			return
		}
		if product.Name == "" && product.Description == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
			return
		}

		product.ID = id
		if err := s.pgHandler.Repository.UpdateProduct(&product); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating product"})
			return
		}

		gatewayHandler, err := handlers.NewGatewayHttpHandler(product.GatewayName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}

		if err := gatewayHandler.UpdateProduct(product.GatewayProductExternalID, product.Name, product.Description, meta.Meta{}); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating product in gateway"})
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

		gatewayHandler, err := handlers.NewGatewayHttpHandler(product.GatewayName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}
		if err := gatewayHandler.DeactivateProduct(product.GatewayProductExternalID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error deactivating product in gateway"})
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
