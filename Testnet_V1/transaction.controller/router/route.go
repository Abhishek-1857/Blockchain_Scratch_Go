package transactionrouter

import (
	"net/http"
	"pop_v1/accounts/accounthelper"
	transactionhelper "pop_v1/transaction.controller/transaction.helper"
)


func TransactionRoute() {
	http.HandleFunc("/recievetransaction", transactionhelper.Recievetransaction)
	http.HandleFunc("/getransactionbyhash", transactionhelper.Gettransaction)
	http.HandleFunc("/sendtesttransaction", transactionhelper.Recievetransaction)
	http.HandleFunc("/getbalancebyaddress", accounthelper.Getbalance)
	http.HandleFunc("/gettransactionbyhash", transactionhelper.Gettransactionbyhash)
	
}
