package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	api_models "github.com/anmolbabu/contact-book/api/models"
	"github.com/anmolbabu/contact-book/dao"
	"github.com/anmolbabu/contact-book/models"
	"github.com/gin-gonic/gin"
)

func validateCreate(c *gin.Context) (httpStatusCode int, contact models.Contact, err error) {
	daoInstance := dao.GetDAOInstance()

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

	var fetchedContact models.Contact
	fetchedContact, err = daoInstance.Get(contact.EmailID)
	if err == nil {
		return http.StatusConflict, fetchedContact, fmt.Errorf("failed to create %+v as %+v already exists", api_models.ToContactResp(contact), api_models.ToContactResp(fetchedContact))
	}

	return http.StatusAccepted, contact, nil
}

func (ch ContactHandler) Add(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()
	httpStatusCode, contact, err := validateCreate(c)
	if err != nil {
		c.JSON(
			httpStatusCode,
			gin.H{"error": err.Error()},
		)
		return
	}
	contact.UpdatedAt = time.Now()
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
