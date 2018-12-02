package api

import (
	"net/http"

	"github.com/anmolbabu/contact-book/api/handlers"
	"github.com/gin-gonic/gin"
)

var contactHandler = handlers.GetContactHandlerInstance()

var routes = map[string]map[string]gin.HandlerFunc{
	"/contacts/:emailid": map[string]gin.HandlerFunc{
		http.MethodGet:    contactHandler.Get,
		http.MethodDelete: contactHandler.Delete,
	},
	"/contacts": map[string]gin.HandlerFunc{
		http.MethodPost: contactHandler.Add,
	},
}
