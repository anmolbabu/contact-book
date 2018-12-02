package db

import (
	"sync"

	"github.com/anmolbabu/contact-book/config"
	"github.com/anmolbabu/contact-book/models"
	"github.com/dgraph-io/badger"
)

type BadgerDB struct {
	conn *badger.DB
}

var badgerDBConn BadgerDB
var once sync.Once

const (
	DB_NAME = "badger"
)

func GetInstance(config *config.Config) (BadgerDB, error) {
	var err error
	once.Do(func() {
		opts := badger.DefaultOptions
		opts.Dir = config.DB_DIR
		opts.ValueDir = config.DB_DIR
		badgerDBConn.conn, err = badger.Open(opts)
	})
	return badgerDBConn, err
}

func (bdb BadgerDB) Add(models.Contact) bool {
	return true
}

func (bdb BadgerDB) GetAll() []models.Contact {
	return []models.Contact{}
}

func (bdb BadgerDB) Get(emailId string) models.Contact {
	return models.Contact{}
}
