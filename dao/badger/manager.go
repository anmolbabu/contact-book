package db

import (
	"encoding/json"
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

func (bdb BadgerDB) Cleanup() error {
	return bdb.conn.Close()
}

func (bdb BadgerDB) Add(contact models.Contact) (err error) {
	contactJSON, err := json.Marshal(contact)
	if err != nil {
		return err
	}
	err = bdb.conn.Update(func(txn *badger.Txn) (err error) {
		txn.Set(
			[]byte(contact.EmailID),
			contactJSON,
		)
		return
	})
	return
}

func (bdb BadgerDB) GetAll(searchContact *models.Contact) (contacts []models.Contact, err error) {
	bdb.conn.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			var contact models.Contact
			item := it.Item()
			err = item.Value(func(val []byte) (err error) {
				err = json.Unmarshal(val, &contact)
				if err != nil {
					return
				}
				return
			})
			if err != nil {
				return
			}
			if contact.IsSearchMatch(searchContact) {
				contacts = append(contacts, contact)
			}
		}
		return
	})
	return
}

func (bdb BadgerDB) Delete(emailId string) (err error) {
	err = bdb.conn.Update(func(txn *badger.Txn) (err error) {
		err = txn.Delete([]byte(emailId))
		return
	})
	return
}

func (bdb BadgerDB) Get(emailId string) (contact models.Contact, err error) {
	err = bdb.conn.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(emailId))
		if err != nil {
			return
		}
		err = item.Value(func(val []byte) error {
			json.Unmarshal(val, &contact)
			return nil
		})
		return
	})
	return
}
