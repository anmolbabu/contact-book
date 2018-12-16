package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/anmolbabu/contact-book/cb_errors"
	"github.com/anmolbabu/contact-book/models"
	"github.com/gin-gonic/gin"
)

func validateUpdate(c *gin.Context) (httpStatusCode int, contact models.Contact, emailId string, err error) {

	emailId = c.Param("emailid")
	fmt.Println(emailId)
	if emailId == "" {
		return http.StatusBadRequest, contact, emailId, fmt.Errorf("invalid emailId")
	}

	err = c.BindJSON(&contact)
	if err != nil {
		return http.StatusBadRequest, contact, emailId, fmt.Errorf("invalid request body. Error: %s", err.Error())
	}

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(emailId) {
		return http.StatusBadRequest, contact, emailId, fmt.Errorf("invalid emailid")
	}

	return http.StatusAccepted, contact, emailId, nil
}

func (ch ContactHandler) Update(c *gin.Context) {

	httpStatusCode, contact, emailId, err := validateUpdate(c)

	if err != nil {
		c.JSON(
			httpStatusCode,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = ch.daoInstance.Update(emailId, contact.Name, contact.EmailID)
	if err != nil {
		if err == cb_errors.CONTACT_NOT_FOUND {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": fmt.Sprintf("failed to update the contact with email id %s. Error: %s", emailId, err.Error())},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed updating contact with email: %s. Error %+v", contact.EmailID, err)},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"Status": "Success"},
	)
	return
}
