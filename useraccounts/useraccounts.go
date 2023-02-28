package useraccounts

import (
	"errors"
	"fmt"

	"github.com/eduard-diumin/banking/helpers"
	"github.com/eduard-diumin/banking/interfaces"
	"github.com/eduard-diumin/banking/transactions"
	"gorm.io/gorm"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount {
	db := helpers.ConnectDB()
	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}
	
	db.Where("id = ?", id).First(&account)
	account.Balance = uint(amount)
	db.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)

	return responseAcc
}

func getAccount(id uint) *interfaces.Account {
	db := helpers.ConnectDB()
	account := &interfaces.Account{}

	dbRresult := db.Where("id = ?", id).First(&account)
	if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return account
}

func Transaction(userID uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIDString := fmt.Sprint(userID)
	isValid := helpers.ValidateToken(userIDString, jwt)

	if isValid {
		fromAccount := getAccount(from)
		toAccount := getAccount(to)
		if fromAccount == nil || toAccount == nil {
			return map[string]interface{}{"message": "Account not found"}
		} else if fromAccount.UserID != userID {
			return map[string]interface{}{"message": "You are not the owner"}
		} else if int(fromAccount.Balance) < amount {
			return map[string]interface{}{"message": "Account balance is to small"}
		}

		// Update account
		updatedAccount := updateAccount(from, int(fromAccount.Balance) - amount)
		updateAccount(to, int(fromAccount.Balance) + amount)

		// Create transaction
		transactions.CreateTransaction(from, to, amount)

		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = updatedAccount
		return response

	} else {
		return map[string]interface{}{"message": "Not valid token"}
	}
}