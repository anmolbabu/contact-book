package handlers

import (
	"fmt"
	"net/http"

	"github.com/anmolbabu/contact-book/cb_errors"

	"github.com/anmolbabu/contact-book/api/models"
	"github.com/gin-gonic/gin"
)

func (ch ContactHandler) Get(c *gin.Context) {
	emailId := c.Param("emailid")
	contact, err := ch.daoInstance.Get(emailId)
	if err != nil {
		if err == cb_errors.CONTACT_NOT_FOUND {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": fmt.Sprintf("requested contact with email id %s does not exist. Error: %+v", emailId, err)},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		models.ToContactResp(contact),
	)
	return
}
