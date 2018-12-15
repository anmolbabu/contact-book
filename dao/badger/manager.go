package badger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/anmolbabu/contact-book/cb_errors"
	"github.com/anmolbabu/contact-book/config"
	"github.com/anmolbabu/contact-book/models"
	"github.com/anmolbabu/contact-book/utils"
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

func (bdb BadgerDB) GetLastKey() (key int) {
	key = 0
	bdb.conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			var err error
			key, err = strconv.Atoi(string(k))
			if err != nil {
				return err
			}
		}

		return nil
	})
	return key
}

func (bdb BadgerDB) Add(contact models.Contact) (err error) {
	contactJSON, err := json.Marshal(contact)
	if err != nil {
		return err
	}

	id := bdb.GetLastKey()
	var eligibleBadgerId []byte
	eligibleBadgerId = []byte(strconv.Itoa(id + 1))

	err = bdb.conn.Update(func(txn *badger.Txn) (err error) {
		err = txn.Set(
			eligibleBadgerId,
			contactJSON,
		)
		return
	})
	return
}

func (bdb BadgerDB) Update(emailId string, newName string, newEmailId string) error {
	key, err := bdb.GetItemKey(emailId)
	if err != nil {
		return err
	}
	foundContact, err := bdb.Get(emailId)
	if err != nil {
		return err
	}

	if newName != "" {
		foundContact.Name = newName
	}

	if newEmailId != "" {
		foundContact.EmailID = newEmailId
	}

	foundContact.UpdatedAt = time.Now()
	contactJSON, err := json.Marshal(foundContact)
	if err != nil {
		return err
	}

	err = bdb.conn.Update(func(txn *badger.Txn) (err error) {
		err = txn.Set(
			key,
			contactJSON,
		)
		return err
	})
	return nil
}

func (bdb BadgerDB) GetAll(searchContact *models.Contact, pageNo int, pageLimit int) (contacts []models.Contact, err error) {
	bdb.conn.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		currInd := 0
		pageBeginInd := 0
		pageEndInd := utils.INVALID_INDEX

		if pageNo != utils.INVALID_INDEX {
			pageBeginInd = pageLimit * (pageNo - 1)
			pageEndInd = pageBeginInd + pageLimit
		}

		for it.Rewind(); it.Valid(); it.Next() {
			for currInd < pageBeginInd {
				currInd++
				if it.Valid() {
					it.Next()
				} else {
					return nil
				}
			}
			if !it.Valid() {
				return fmt.Errorf("Not Available")
			}
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
			if pageEndInd != utils.INVALID_INDEX {
				if currInd == pageEndInd-1 {
					break
				}
			}
			currInd++
		}
		return
	})
	if len(contacts) == 0 {
		err = cb_errors.CONTACT_NOT_FOUND
	}
	return
}

func (bdb BadgerDB) GetItemKey(emailID string) (key []byte, err error) {
	err = bdb.conn.View(func(txn *badger.Txn) (err error) {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		var contact models.Contact

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			tKey := item.Key()
			err = item.Value(func(val []byte) (err error) {
				err = json.Unmarshal(val, &contact)
				if err != nil {
					return
				}
				defContact := models.GetDefaultContact()
				defContact.EmailID = emailID
				if contact.IsSearchMatch(defContact) {
					key = tKey
					return
				}
				return
			})

			if len(key) > 0 {
				break
			}
		}
		return
	})
	if len(key) == 0 {
		return key, cb_errors.CONTACT_NOT_FOUND
	}
	return
}

func (bdb BadgerDB) Delete(emailId string) (err error) {
	key, err := bdb.GetItemKey(emailId)
	if err != nil {
		return
	}
	err = bdb.conn.Update(func(txn *badger.Txn) (err error) {
		err = txn.Delete(key)
		return
	})
	return
}

func (bdb BadgerDB) Get(emailId string) (contact models.Contact, err error) {
	key, err := bdb.GetItemKey(emailId)
	if err != nil {
		return
	}
	err = bdb.conn.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get(key)
		if err != nil {
			return
		}
		err = item.Value(func(val []byte) error {
			json.Unmarshal(val, &contact)
			return nil
		})
		return
	})
	if contact.IsEmpty() {
		err = cb_errors.CONTACT_NOT_FOUND
	}
	return
}
