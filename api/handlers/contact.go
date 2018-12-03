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
	contact := models.Contact{}
	contactPtr := &contact

	err := c.Bind(&contact)
	if err != nil {
		contactPtr = nil
	}
	contacts, err := daoInstance.GetAll(contactPtr)
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

	contact.UpdatedAt = time.Now()

	if contact.EmailID == "" {
		contact.EmailID = emailId
	} else {
		if emailId != contact.EmailID {
			c.JSON(
				500,
				gin.H{"error": "email ids in the request body and url param are different"},
			)
			return
		}
	}

	oldContact, err := daoInstance.Get(emailId)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": fmt.Sprintf("failed fetching contact with email id %s", emailId)},
		)
		return
	}
	contact = SuperImposeContacts(contact, oldContact)

	err = daoInstance.Add(contact)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
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
