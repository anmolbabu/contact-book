package api

import (
	"net/http"

	"github.com/anmolbabu/contact-book/api/handlers"
	"github.com/gin-gonic/gin"
)

func getRouter(handler handlers.ContactHandler) map[string]map[string]gin.HandlerFunc {
	return map[string]map[string]gin.HandlerFunc{
		"/contacts/:emailid": map[string]gin.HandlerFunc{
			http.MethodGet:    handler.Get,
			http.MethodDelete: handler.Delete,
			http.MethodPut:    handler.Update,
		},
		"/contacts": map[string]gin.HandlerFunc{
			http.MethodPost: handler.Add,
			http.MethodGet:  handler.GetAll,
		},
	}
}
