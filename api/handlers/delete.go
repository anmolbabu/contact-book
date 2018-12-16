package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/anmolbabu/contact-book/cb_errors"
	"github.com/gin-gonic/gin"
)

func validateDelete(c *gin.Context) (httpStatusCode int, emailId string, err error) {
	emailId = c.Param("emailid")
	if emailId == "" {
		return http.StatusBadRequest, emailId, fmt.Errorf("invalid emailId")
	}

	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !rxEmail.MatchString(emailId) {
		return http.StatusBadRequest, emailId, fmt.Errorf("invalid emailid")
	}

	return http.StatusAccepted, emailId, nil
}

func (ch ContactHandler) Delete(c *gin.Context) {

	httpStatusCode, emailId, err := validateDelete(c)
	if err != nil {
		c.JSON(
			httpStatusCode,
			gin.H{"error": err.Error()},
		)
		return
	}

	err = ch.daoInstance.Delete(emailId)
	if err != nil {
		if err == cb_errors.CONTACT_NOT_FOUND {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": fmt.Sprintf("contact with emailid %s not found", emailId)},
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
		200,
		gin.H{"Status": "Success"},
	)
	return
}
