package dao

import (
	"fmt"

	"github.com/anmolbabu/contact-book/config"
	badger_db "github.com/anmolbabu/contact-book/dao/badger"
	"github.com/anmolbabu/contact-book/models"
)

type DataAccess interface {
	GetAll() []models.Contact
	Get(emailId string) models.Contact
	Add(models.Contact) bool
}

var db DataAccess

func Init(config *config.Config) (DataAccess, error) {
	var err error
	switch config.DB {
	case badger_db.DB_NAME:
		db, err = badger_db.GetInstance(config)
		return db, err
	default:
		return db, fmt.Errorf("failed to initialise db")
	}
}

func GetDAOInstance() DataAccess {
	return db
}
