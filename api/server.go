package api

import (
	db "github.com/GustavoNoronha0/gofinance-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.SQLStore
	router *gin.Engine
}

func NewServer(store *db.SQLStore) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.GET("/user/:username", server.getUser)
	router.GET("/user/id/:id", server.getUserById)

	router.POST("/category", server.createCategory)
	router.GET("/category/id/:id", server.getCategory)
	router.GET("/category", server.getCategories)
	router.DELETE("/category/:id", server.deleteCategory)
	router.PUT("/category/:id", server.updateCategory)

	router.POST("/account", server.createAccount)
	router.GET("/account/id/:id", server.getAccount)
	router.GET("/account", server.getAccounts)
	router.GET("/account/graph/:user_id/:type", server.getAccountGraph)
	router.GET("/account/reports/:user_id/:type", server.getAccountReports)
	router.DELETE("/account/:id", server.deleteAccount)
	router.PUT("/account/:id", server.updateAccount)

	router.POST("/login", server.login)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"api has error:": err.Error()}
}
