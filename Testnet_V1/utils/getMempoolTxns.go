package utils

import (
	"encoding/json"
	"log"
	"pop_v1/database"
	"pop_v1/models"
	"sync"

	"github.com/tecbot/gorocksdb"
)

func GetMempoolTxns() (<-chan func() (string, models.Transaction), <-chan bool, *sync.WaitGroup) {
	out := make(chan func() (string, models.Transaction))
	done := make(chan bool)
	wg := sync.WaitGroup{}
	go func() {
		// Reading data
		readOpts := gorocksdb.NewDefaultReadOptions()
		defer readOpts.Destroy()
		iter := database.Mempool_db.NewIterator(readOpts)
		defer iter.Close()
		cnt := 0
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := iter.Key()
			value := iter.Value()
			TransactionTest := models.Transaction{}
			if err := json.Unmarshal(value.Data(), &TransactionTest); err != nil {
				log.Fatalf("Error deserializing data %v\n", err)
				done <- false
				close(done)
				close(out)
				return
			}
			wg.Add(1)
			out <- (func() (string, models.Transaction) {
				return string(key.Data()), TransactionTest
			})
			wg.Wait()
			cnt++
			if cnt == 20 {
				break
			}
		}
		done <- true
		close(done)
		close(out)
	}()
	return out, done, &wg
}
