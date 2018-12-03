package models

import "time"

// Contact is expected to hold contact details like name, emailid and timestamp at which the contact was created/modified
type Contact struct {
	Name      string    `json:"Name"`
	EmailID   string    `json:"EmailID"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func (c Contact) IsSearchMatch(searchPtr *Contact) bool {
	isMatch := true
	if searchPtr == nil {
		return isMatch
	}

	if searchPtr.Name != "" {
		isMatch = isMatch && (searchPtr.Name == c.Name)
	}
	if searchPtr.EmailID != "" {
		isMatch = isMatch && (searchPtr.EmailID == c.EmailID)
	}

	return isMatch
}
