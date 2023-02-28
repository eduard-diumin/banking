package main

import (
	"github.com/eduard-diumin/banking/api"
)

func main() {
	// Do migrations
	// migrations.Migrate()
	// migrations.MigrateTransactions()
	api.StartApi()
}