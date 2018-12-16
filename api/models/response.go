package models

import "github.com/anmolbabu/contact-book/models"

type ContactResp struct {
	Name    string `json:"name"`
	EmailId string `json:"emailid"`
}

func ToContactResp(contact models.Contact) (cr ContactResp) {
	cr.Name = contact.Name
	cr.EmailId = contact.EmailID
	return
}

func ToContactResps(contacts []models.Contact) (ctrs []ContactResp) {
	for _, contact := range contacts {
		ctrs = append(
			ctrs,
			ContactResp{
				Name:    contact.Name,
				EmailId: contact.EmailID,
			},
		)
	}
	return
}
