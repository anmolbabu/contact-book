package handlers

import (
	"sync"
	"time"

	"github.com/anmolbabu/contact-book/dao"
	"github.com/anmolbabu/contact-book/models"
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
	emailId := c.Param("emailid")
	contact, err := daoInstance.Get(emailId)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		200,
		gin.H{"contact": contact},
	)
	return
}

func (ch ContactHandler) Add(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	contact := models.Contact{}
	c.BindJSON(&contact)
	contact.UpdatedAt = time.Now()
	err := daoInstance.Add(contact)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
	}
	c.JSON(
		200,
		gin.H{"Status": "Success"},
	)
	return
}

func (ch ContactHandler) Delete(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	emailId := c.Param("emailid")
	err := daoInstance.Delete(emailId)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
	}
	c.JSON(
		200,
		gin.H{"Status": "Success"},
	)
	return
}
