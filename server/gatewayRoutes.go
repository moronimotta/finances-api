package server

import (
	"finances-api/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) initGatewayRoutes() {

	s.app.POST("/gateways", func(ctx *gin.Context) {
		gateway := entities.Gateway{}

		if err := ctx.ShouldBindJSON(&gateway); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := s.pgHandler.Repository.CreateGateway(gateway); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating gateway"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway created successfully"})
	})

	s.app.GET("/gateways/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		gateway, err := s.pgHandler.Repository.GetGatewayByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Gateway not found"})
			return
		}
		ctx.JSON(http.StatusOK, gateway)
	})

	s.app.GET("/gateways", func(ctx *gin.Context) {
		gateways, err := s.pgHandler.Repository.GetAllGateways()
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		gateway.ID = id // Ensure the ID is set for the update

		if err := s.pgHandler.Repository.UpdateGateway(&gateway); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating gateway"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway updated successfully"})
	})
	s.app.DELETE("/gateways/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := s.pgHandler.Repository.DeleteGateway(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting gateway"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Gateway deleted successfully"})
	})

}
