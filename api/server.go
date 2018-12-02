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

	for route, apiHandlerMap := range routes {
		for reqType, apiHandler := range apiHandlerMap {
			switch reqType {
			case http.MethodGet:
				router.GET(route, apiHandler)
			}
		}
	}
	router.Run(":" + config.Port)
}
