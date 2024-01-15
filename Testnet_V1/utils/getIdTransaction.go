package utils

import (
	"pop_v1/database"
	"sync"

	"github.com/tecbot/gorocksdb"
)

var lock_ sync.Mutex

func getCurrentId_() string {

	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Initialize an iterator
	iter := database.TestTransaction_db.NewIterator(readOpts)
	defer iter.Close()

	// Seek to the last key
	iter.SeekToLast()

	// Check if the iterator is valid
	if iter.Valid() {
		// Get the key
		key := iter.Key()
		return string(key.Data())
	} else {
		//if there is no record in the database
		return "abcdefg"
	}
}

func GenerateIdClient() string {
	id := getCurrentId_()
	lock_.Lock()
	id = next_id(id)
	lock_.Unlock()
	return id
}
