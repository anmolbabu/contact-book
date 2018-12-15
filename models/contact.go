package models

import (
	"time"

	"github.com/anmolbabu/contact-book/utils"
)

// Contact is expected to hold contact details like name, emailid and timestamp at which the contact was created/modified
type Contact struct {
	Name      string    `json:"name"`
	EmailID   string    `json:"emailid"`
	UpdatedAt time.Time `json:"updatedat"`
}

func GetDefaultContact() *Contact {
	return &Contact{
		Name:    utils.INVALID_STRING,
		EmailID: utils.INVALID_STRING,
	}
}

func (c Contact) IsSearchMatch(searchPtr *Contact) bool {
	isMatch := true
	if searchPtr == nil {
		return isMatch
	}

	if searchPtr.Name != utils.INVALID_STRING {
		isMatch = isMatch && (searchPtr.Name == c.Name)
	}
	if searchPtr.EmailID != utils.INVALID_STRING {
		isMatch = isMatch && (searchPtr.EmailID == c.EmailID)
	}

	return isMatch
}

func (c Contact) IsSame(c1 Contact) bool {
	return (c.EmailID == c1.EmailID) && (c.Name == c1.Name)
}

func (c Contact) IsEmpty() bool {
	return c.Name == "" && c.EmailID == ""
}
