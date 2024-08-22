package api

import (
	"github.com/Ayobami6/pickitup_v3/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIServer struct {
	address string
	db      *gorm.DB
}


func NewAPIServer(address string, db *gorm.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

// Implement the Run method to start the server

func (s *APIServer) Run() error {
    // TODO: Implement server logic here
	router := gin.Default()
	routes.RootRoute(router)
	return router.Run(s.address)
}