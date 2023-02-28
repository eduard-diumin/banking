package transactions

import (
	"github.com/eduard-diumin/banking/helpers"
	"github.com/eduard-diumin/banking/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	db := helpers.ConnectDB()
	transactions := &interfaces.Transactions{From: From, To: To, Amount: Amount}
	db.Create(&transactions)
}