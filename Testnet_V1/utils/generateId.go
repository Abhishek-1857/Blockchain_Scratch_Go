package utils

import (
	"pop_v1/database"
	"sync"

	"github.com/tecbot/gorocksdb"
)

var lock sync.Mutex

func getCurrentId() string {
	

	//read options
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Initialize an iterator
	iter := database.Blockchain_db.NewIterator(readOpts)
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
		return "abcdefgh"
	}
}

func next_id(input string) string {
	arr := []byte(input)
	n := len(arr)
	i := n - 2
	// Find the first element that is not in ascending order
	for i >= 0 && arr[i] >= arr[i+1] {
		i--
	}
	if i == -1 {
		// This case will never reach and if yes we will return a random string
		// with greater length
		return "abcdefghij"
	}
	j := n - 1
	// Find the rightmost element greater than the pivot
	for arr[j] <= arr[i] {
		j--
	}
	// Swap the pivot and the rightmost element
	arr[i], arr[j] = arr[j], arr[i]
	// Reverse the elements to the right of the pivot
	i++
	j = n - 1
	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
	return string(arr)
}

func GenerateId() string {
	id := getCurrentId()
	lock.Lock()
	id = next_id(id)
	lock.Unlock()
	return id
}
