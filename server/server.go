package server

import (
	"finances-api/handlers"
	logs "finances-api/utils/logs"
	"io/ioutil"
	"net/http"

	"finances-api/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app *gin.Engine
	db  db.Database
}

func NewServer(db db.Database) *Server {
	logs.InitLogging()

	return &Server{
		app: gin.Default(),
		db:  db,
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

	stripeHandler, err := handlers.NewHttpHandler("stripe")
	if err != nil {
		return
	}

	s.app.POST("/webhook/stripe", func(ctx *gin.Context) {
		const MaxBodyBytes = int64(65536)

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
		payload, err := ioutil.ReadAll(ctx.Request.Body)

		stripeHandler.EventBus(payload, ctx.Request.Header.Get("Stripe-Signature"))

		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
			return
		}
	})

}
