package handlers

import (
	"fmt"
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
		contact,
	)
	return
}

func (ch ContactHandler) GetAll(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()

	lr := GetDefaultListReq()
	lr, err := lr.Serialise(c.Request.URL.Query())
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}

	contacts, err := daoInstance.GetAll(&(lr.Contact), lr.PageNo, lr.PageLimit)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		200,
		contacts,
	)
	return
}

func (ch ContactHandler) Add(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	contact := models.Contact{}
	err := c.BindJSON(&contact)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
	}
	contact.UpdatedAt = time.Now()
	err = daoInstance.Add(contact)
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

func (ch ContactHandler) Update(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	emailId := c.Param("emailid")
	contact := models.Contact{}

	err := c.BindJSON(&contact)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": "malformed json"},
		)
		return
	}
	if contact.EmailID != "" {
		c.JSON(
			500,
			gin.H{"error": "updating emailid is not supported"},
		)
		return
	}

	contact.UpdatedAt = time.Now()

	err = daoInstance.Update(emailId, contact.Name)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": fmt.Sprintf("failed updating contact with email: %s. Error %+v", emailId, err)},
		)
	}
	c.JSON(
		200,
		gin.H{"Status": "Success"},
	)
	return
}

func SuperImposeContacts(newContact models.Contact, oldContact models.Contact) models.Contact {
	if newContact.Name == "" {
		newContact.Name = oldContact.Name
	}
	if newContact.EmailID == "" {
		newContact.EmailID = oldContact.EmailID
	}
	newContact.UpdatedAt = time.Now()
	return newContact
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
