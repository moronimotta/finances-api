package server

import (
	"finances-api/entities"
	"finances-api/handlers"
	"finances-api/usecases"
	logs "finances-api/utils/logs"
	"io/ioutil"
	"net/http"

	"finances-api/db"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	app         *gin.Engine
	db          db.Database
	usecases    *usecases.PaymentAPIUsecases
	redisClient *redis.Client
}

func NewServer(db db.Database, redisClient *redis.Client) *Server {
	logs.InitLogging()

	usecases := usecases.NewPaymentAPIUsecases("stripe", db)

	return &Server{
		app:         gin.Default(),
		db:          db,
		redisClient: redisClient,
		usecases:    usecases,
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
	s.initInvoicesRoutes()

	s.app.POST("/checkout", func(ctx *gin.Context) {
		checkout := entities.Checkout{}
		if err := ctx.ShouldBindJSON(&checkout); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		checkoutURL, err := s.usecases.Gateway.CreateCheckoutSession(checkout.PriceID, checkout.CustomerID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error creating checkout session"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"client_secret": checkoutURL})
	})

	s.app.POST("/webhook/stripe", func(ctx *gin.Context) {
		const MaxBodyBytes = int64(65536)

		// calls the webhookhandler
		stripeHandler, err := handlers.NewWebhookHandler("stripe", s.usecases)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
		payload, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
			return
		}

		if err := stripeHandler.EventBus(payload, ctx.Request.Header.Get("Stripe-Signature")); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	})

}
