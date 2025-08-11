package server

import (
	"finances-api/entities"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initGatewayRoutes() {

	s.app.POST("/gateways", func(ctx *gin.Context) {
		gateway := entities.Gateway{}

		if err := ctx.ShouldBindJSON(&gateway); err != nil {
			slog.Error("Invalid request when creating gateway", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.usecases.Db.CreateGateway(gateway); err != nil {
			slog.Error("Error creating gateway", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating gateway"})
			return
		}

		slog.Info("Gateway created successfully")
		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway created successfully"})
	})

	s.app.GET("/gateways/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		gateway, err := s.usecases.Db.GetGatewayByID(id)
		if err != nil {
			slog.Error("Gateway not found", "id", id, "error", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		slog.Info("Gateway retrieved successfully")
		ctx.JSON(http.StatusOK, gateway)
	})

	s.app.GET("/gateways", func(ctx *gin.Context) {
		gateways, err := s.usecases.Db.GetAllGateways()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting gateways"})
			return
		}
		ctx.JSON(http.StatusOK, gateways)
	})

	s.app.PUT("/gateways/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		gateway := entities.Gateway{}

		if err := ctx.ShouldBindJSON(&gateway); err != nil {
			slog.Error("Invalid request when updating gateway", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		gateway.ID = id // Ensure the ID is set for the update

		if err := s.usecases.Db.UpdateGateway(&gateway); err != nil {
			slog.Error("Error updating gateway", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating gateway"})
			return
		}

		slog.Info("Gateway updated successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway updated successfully"})
	})
	s.app.DELETE("/gateways/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.usecases.Db.DeleteGateway(id); err != nil {
			slog.Error("Error deleting gateway", "id", id, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting gateway"})
			return
		}
		slog.Info("Gateway deleted successfully", "id", id)
		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway deleted successfully"})
	})

}
