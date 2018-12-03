package dao

import (
	"fmt"

	"github.com/anmolbabu/contact-book/config"
	badger_db "github.com/anmolbabu/contact-book/dao/badger"
	"github.com/anmolbabu/contact-book/models"
)

type DataAccess interface {
	GetAll(searchContact *models.Contact) ([]models.Contact, error)
	Get(emailId string) (models.Contact, error)
	Add(models.Contact) error
	Delete(emailId string) error
	Cleanup() error
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

func Cleanup() error {
	return db.Cleanup()
}
