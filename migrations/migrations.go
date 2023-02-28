package migrations

import (
	"github.com/eduard-diumin/banking/helpers"
	"github.com/eduard-diumin/banking/interfaces"
)




func createAccounts()  {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{
			FirstName: "Jack",
			LastName: "Black",
			Email: "jackblack@gmail.com",
		}, {
			FirstName: "Jack",
			LastName: "White",
			Email: "jackwhite@gmail.com",
		},
	}

	for i := 0; i <len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].FirstName))
		user := &interfaces.User{
			FirstName: users[i].FirstName,
			LastName: users[i].LastName,
			Email: users[i].Email,
			Password: generatedPassword,
		}
		db.Create(&user)

		account := &interfaces.Account{
			Type: "User Account",
			Name: string(users[i].FirstName + users[i].LastName + "Account"), 
			Balance: uint(10000 * int(i+1)),
			UserID: user.ID,
		}
		db.Create(&account)
	}
}

func Migrate()  {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	db := helpers.ConnectDB()
	db.AutoMigrate(&User, &Account)

	createAccounts()
}

func MigrateTransactions()  {
	Transactions := &interfaces.Transactions{}

	db := helpers.ConnectDB()
	db.AutoMigrate(&Transactions)
}