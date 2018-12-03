package api

import (
	"net/http"

	"github.com/anmolbabu/contact-book/config"
	"github.com/gin-gonic/gin"
	//"github.com/vsouza/go-gin-boilerplate/middlewares"
)

type APIRequestHandler interface {
	Get(c *gin.Context)
}

func NewRouter(config *config.Config) {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		config.AuthUser: config.AuthPass,
	}))

	for route, apiHandlerMap := range routes {
		for reqType, apiHandler := range apiHandlerMap {
			switch reqType {
			case http.MethodGet:
				authorized.GET(route, apiHandler)
			case http.MethodPost:
				authorized.POST(route, apiHandler)
			case http.MethodDelete:
				authorized.DELETE(route, apiHandler)
			case http.MethodPut:
				authorized.PUT(route, apiHandler)
			}
		}
	}
	router.Run(":" + config.Port)
}
