package database

import "github.com/tecbot/gorocksdb"

var Account_db *gorocksdb.DB
var Blockchain_db *gorocksdb.DB
var Mempool_db *gorocksdb.DB
var TestTransaction_db *gorocksdb.DB
