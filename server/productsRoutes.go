package server

import (
	"finances-api/entities"
	"finances-api/utils/meta"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initProductsRoutes() {
	s.app.POST("/products", func(ctx *gin.Context) {
		product := entities.Products{}

		if err := ctx.ShouldBindJSON(&product); err != nil {
			slog.Error("Failed to bind JSON", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var err error
		product.GatewayProductExternalID, err = s.usecases.Gateway.CreateProduct(product.Name, product.Description, product)
		if err != nil {
			slog.Error("Failed to create product in gateway", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in gateway"})
			return
		}

		product.GatewayPriceExternalID, err = s.usecases.Gateway.CreatePrice(product.GatewayProductExternalID, product.Price, product.Currency)
		if err != nil {
			slog.Error("Failed to create price in gateway", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating price in gateway"})
			return
		}

		err = s.usecases.Db.CreateProduct(&product)
		if err != nil {
			slog.Error("Failed to create product in database", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating product in database"})
			return
		}

		meta := meta.New()
		meta.Add("local_product_id", product.ID)
		meta.Add("gateway_current_price_id", product.GatewayPriceExternalID)

		s.usecases.Gateway.UpdateProduct(product.GatewayProductExternalID, "", "", meta)

		slog.Info("Product created successfully", "product_id", product.ID, "external_id", product.GatewayProductExternalID)
		ctx.JSON(http.StatusOK, gin.H{"message": "Product created successfully: ", "product_id": product.ID, "external_id": product.GatewayProductExternalID})
	})

	s.app.GET("/products", func(ctx *gin.Context) {
		products, err := s.usecases.Db.GetAllProducts()
		if err != nil {
			slog.Error("Failed to get all products", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting products"})
			return
		}
		slog.Info("Retrieved all products", "count", len(products))
		ctx.JSON(http.StatusOK, gin.H{"products": products})
	})

	s.app.GET("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		product, err := s.usecases.Db.GetProductByID(id)
		if err != nil {
			slog.Error("Failed to get product by ID", "error", err, "id", id)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		slog.Info("Retrieved product by ID", "id", id)
		ctx.JSON(http.StatusOK, product)
	})

	s.app.PUT("/products/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")
		product := entities.Products{}
		if err := ctx.ShouldBindJSON(&product); err != nil {
			slog.Error("Failed to bind JSON", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if product.ID == "" {
			product.ID = id
		}
		if id != product.ID {
			slog.Error("Product ID mismatch", "url_id", id, "body_id", product.ID)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product ID mismatch"})
			return
		}

		if product.Price != 0 {
			newPriceID, err := s.usecases.Gateway.ChangePrice(product.GatewayPriceExternalID, product.GatewayProductExternalID, product.Price, product.Currency)
			if err != nil {
				slog.Error("Failed to change price in gateway", "error", err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error changing price in gateway"})
				return
			}
			product.GatewayPriceExternalID = newPriceID
		}

		product.ID = id
		if err := s.usecases.Db.UpdateProduct(&product); err != nil {
			slog.Error("Failed to update product in database", "error", err, "id", id)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating product"})
			return
		}

		meta := meta.New()
		meta.Add("gateway_current_price_id", product.GatewayPriceExternalID)

		if err := s.usecases.Gateway.UpdateProduct(product.GatewayProductExternalID, product.Name, product.Description, meta); err != nil {
			slog.Error("Failed to update product in gateway", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error updating product in gateway"})
			return
		}

		slog.Info("Product updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	})

	s.app.DELETE("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		product, err := s.usecases.Db.GetProductByID(id)
		if err != nil {
			slog.Error("Failed to get product for deletion", "error", err, "id", id)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		if err := s.usecases.Gateway.DeactivateProduct(product.GatewayProductExternalID); err != nil {
			slog.Error("Failed to deactivate product in gateway", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error deactivating product in gateway"})
			return
		}
		if err := s.usecases.Db.DeleteProduct(product.ID); err != nil {
			slog.Error("Failed to delete product from database", "error", err, "id", product.ID)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting product"})
			return
		}
		slog.Info("Product deleted successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})

	s.app.GET("/products/g/:external_id", func(ctx *gin.Context) {
		externalID := ctx.Param("external_id")
		product, err := s.usecases.Db.GetProductByExternalID(externalID)
		if err != nil {
			slog.Error("Failed to get product by external ID", "error", err, "external_id", externalID)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error getting product"})
			return
		}
		slog.Info("Retrieved product by external ID", "external_id", externalID)
		ctx.JSON(http.StatusOK, product)
	})
}
