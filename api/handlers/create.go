package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/anmolbabu/contact-book/cb_errors"

	"github.com/anmolbabu/contact-book/models"
	"github.com/gin-gonic/gin"
)

func validateCreate(c *gin.Context) (httpStatusCode int, contact models.Contact, err error) {
	err = c.BindJSON(&contact)
	if err != nil {
		return http.StatusBadRequest, contact, fmt.Errorf("invalid parameter passed")
	}

	if contact.Name == "" || contact.EmailID == "" {
		return http.StatusBadRequest, contact, fmt.Errorf("expecting the request to pass name and emailid")
	}
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(contact.EmailID) {
		return http.StatusBadRequest, contact, fmt.Errorf("invalid emailid")
	}

	return http.StatusAccepted, contact, nil
}

func (ch ContactHandler) Add(c *gin.Context) {
	httpStatusCode, contact, err := validateCreate(c)
	if err != nil {
		c.JSON(
			httpStatusCode,
			gin.H{"error": err.Error()},
		)
		return
	}
	contact.UpdatedAt = time.Now()
	err = ch.daoInstance.Add(contact)
	if err != nil {
		if err == cb_errors.DUPLICATE_CONTACT {
			c.JSON(
				http.StatusConflict,
				gin.H{"error": err.Error()},
			)
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusCreated,
		gin.H{"Status": "Success"},
	)
	return
}
