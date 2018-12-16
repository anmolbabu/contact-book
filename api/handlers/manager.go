package handlers

import "github.com/anmolbabu/contact-book/dao"

type ContactHandler struct {
	daoInstance dao.DataAccess
}

var contactHandler *ContactHandler

func Init(da dao.DataAccess) *ContactHandler {
	contactHandler = &ContactHandler{
		daoInstance: da,
	}
	return contactHandler
}
