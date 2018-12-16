package dao

import (
	"fmt"

	"github.com/anmolbabu/contact-book/config"
	badger_db "github.com/anmolbabu/contact-book/dao/badger"
	"github.com/anmolbabu/contact-book/models"
)

type DataAccess interface {
	GetAll(searchContact *models.Contact, pageNo int, pageLimit int) ([]models.Contact, error)
	GetItemKey(emailID string) (key []byte, err error)
	Get(emailId string) (models.Contact, error)
	Add(models.Contact) error
	Update(string, string, string) error
	Delete(emailId string) error
	Cleanup() error
}

type DataAccessLayer struct {
	db DataAccess
}

var dataAccessManager DataAccessLayer

func Init(config *config.Config) (DataAccess, error) {
	var err error
	switch config.DB {
	case badger_db.DB_NAME:
		dataAccessManager.db, err = badger_db.GetInstance(config)
		return dataAccessManager.db, err
	default:
		return dataAccessManager.db, fmt.Errorf("failed to initialise db")
	}
}

func SetDAO(dao DataAccess) {
	dataAccessManager.db = dao
}

func GetDAOInstance() DataAccess {
	return dataAccessManager.db
}

func Cleanup() error {
	return dataAccessManager.db.Cleanup()
}
