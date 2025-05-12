package server

import (
	"finances-api/utils"

	"finances-api/db"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app *gin.Engine
	db  db.Database
}

func NewServer(db db.Database) *Server {
	utils.InitLogging()

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
	// Initialize payment handler here
}
