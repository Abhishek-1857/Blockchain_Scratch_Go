package router

import (
	clientrouter "pop_v1/client.controller/router"
	noderouter "pop_v1/node.controller/router"
	TransactionRoute "pop_v1/transaction.controller/router"
)

// SetupRoutes func
func MainRoute() {
	TransactionRoute.TransactionRoute()
	clientrouter.SetupClientRoute()
	noderouter.NodeRoutes()
}
