package handlers

import (
	"sync"

	"github.com/anmolbabu/contact-book/dao"
	"github.com/gin-gonic/gin"
)

type ContactHandler struct{}

var contactHandler *ContactHandler
var once sync.Once

func GetContactHandlerInstance() *ContactHandler {
	once.Do(func() {
		contactHandler = &ContactHandler{}
	})
	return contactHandler
}

func (ch ContactHandler) Get(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	contact, err := daoInstance.Get("abc@email.com")
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
	}
	c.JSON(
		200,
		gin.H{"contact": contact},
	)
	return
}
