package report

import (
	"github.com/dgraph-io/badger/v2"
	"log"
)

func Put(k string, v string) error {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Printf("Error opening database %s", err.Error())
		return err
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(k), []byte(v))
		err = txn.Commit()
		if err != nil {
			log.Printf("Error writing result to db %s", err)
			return err
		}
		// Your code here…
		return nil
	})
	if err != nil {
		log.Printf("Key %s value %s successfully stored in db", k, v)
		return err
	}
	return nil
}

func Get(k string) (string, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Printf("Error opening database %s", err.Error())
		return "", err
	}
	defer db.Close()
	var value []byte

	err = db.View(func(txn *badger.Txn) error {
		item, keyNotFoundError := txn.Get([]byte(k))
		if keyNotFoundError != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = val
			return nil
		})
		return nil
	})

	if err != nil {
		log.Printf("Error retrieving value for key %s %s", k, err.Error())
		return "", err
	}
	return string(value), nil
}
