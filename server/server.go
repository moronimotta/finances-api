package server

import (
	"finances-api/entities"
	"finances-api/handlers"
	logs "finances-api/utils/logs"
	"io/ioutil"
	"net/http"

	"finances-api/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app       *gin.Engine
	pgHandler *handlers.DbHttpHandler
}

func NewServer(db db.Database) *Server {
	logs.InitLogging()

	pgHandler, err := handlers.NewDbHttpHandler("postgres", db)
	if err != nil {
		return nil
	}

	return &Server{
		app:       gin.Default(),
		pgHandler: pgHandler,
	}
}
func (s *Server) Start() {

	s.initializeMiddlewares()

	s.initializePaymentHttpHandler()

	if err := s.app.Run(":9090"); err != nil {
		panic(err)
	}
}

func (s *Server) initializePaymentHttpHandler() {

	s.initProductsRoutes()
	s.initUserProductsRoutes()
	s.initTransactionsRoutes()
	s.initGatewayRoutes()

	s.app.POST("/checkout", func(ctx *gin.Context) {
		checkout := entities.Checkout{}
		if err := ctx.ShouldBindJSON(&checkout); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		gatewayHandler, err := handlers.NewGatewayHttpHandler(checkout.GatewayName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}

		checkoutURL, err := gatewayHandler.Repository.CreateCheckoutSession(checkout.PriceID, checkout.CustomerID, checkout.SuccessURL, checkout.CancelURL)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating checkout session"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"checkout_url": checkoutURL})
	})

	s.app.POST("/webhook/stripe", func(ctx *gin.Context) {
		const MaxBodyBytes = int64(65536)

		stripeHandler, err := handlers.NewGatewayHttpHandler("stripe")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
		payload, err := ioutil.ReadAll(ctx.Request.Body)

		stripeHandler.EventBus(payload, ctx.Request.Header.Get("Stripe-Signature"))

		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
			return
		}
	})

}
