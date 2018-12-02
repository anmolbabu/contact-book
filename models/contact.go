package models

import "time"

// Contact is expected to hold contact details like name, emailid and timestamp at which the contact was created/modified
type Contact struct {
	Name      string    `json:"Name"`
	EmailID   string    `json:"EmailID"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
