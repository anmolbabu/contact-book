package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/anmolbabu/contact-book/cb_errors"
	"github.com/anmolbabu/contact-book/dao"
	"github.com/gin-gonic/gin"
)

func validateDelete(c *gin.Context) (httpStatusCode int, emailId string, err error) {
	daoInstance := dao.GetDAOInstance()

	emailId = c.Param("emailid")
	if emailId == "" {
		return http.StatusBadRequest, emailId, fmt.Errorf("invalid emailId")
	}

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(emailId) {
		return http.StatusBadRequest, emailId, fmt.Errorf("invalid emailid")
	}

	_, err = daoInstance.Get(emailId)
	if err != nil {
		if err == cb_errors.CONTACT_NOT_FOUND {
			return http.StatusNotFound, emailId, fmt.Errorf("failed to delete the requested contact %s. Error: %s", emailId, err.Error())
		} else {
			return http.StatusInternalServerError, emailId, fmt.Errorf("failed to delete the requested contact %s. Error %s", emailId, err.Error())
		}
	}

	return http.StatusAccepted, emailId, nil
}

func (ch ContactHandler) Delete(c *gin.Context) {
	daoInstance := dao.GetDAOInstance()

	httpStatusCode, emailId, err := validateDelete(c)
	if err != nil {
		c.JSON(
			httpStatusCode,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = daoInstance.Delete(emailId)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
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
