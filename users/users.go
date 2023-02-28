package users

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/eduard-diumin/banking/helpers"
	"github.com/eduard-diumin/banking/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func prepareToken(user *interfaces.User) string {
		// Sign token
		tokenContent := jwt.MapClaims{
			"user_id": user.ID,
			"expiry": time.Now().Add(time.Minute ^ 60).Unix(),
		}
		jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
		token, err := jwtToken.SignedString([]byte("TokenPassword"))
		helpers.HandleErr(err)

		return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{}  {
	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Account: accounts,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		var token = prepareToken(user) 
		response["jwt"] = token
	}

	response["data"] = responseUser

	return response
}

func Login(username string, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect db
		db := helpers.ConnectDB()
		user := &interfaces.User{}

		dbRresult := db.Where("username = ?", username).First(&user)
		if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found"}
		}

		// Check password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr !=nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		// Find account for the user
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)


		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}

func Register(username string, email string, pass string) map[string]interface{}  {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	
		if valid {
			db := helpers.ConnectDB()
			generatedPassword := helpers.HashAndSalt([]byte(username))
			user := &interfaces.User{
				FirstName: username,
				Email: email,
				Password: generatedPassword,
			}
			db.Create(&user)
	
			account := &interfaces.Account{
				Type: "User Account",
				Name: string(username + "Account"), 
				Balance: 0,
				UserID: user.ID,
			}
			db.Create(&account)

			accounts := []interfaces.ResponseAccount{}
			respAccount := interfaces.ResponseAccount{
				ID: account.ID, 
				Name: account.Name, 
				Balance: int(account.Balance),
			}
			accounts = append(accounts, respAccount)
			var response = prepareResponse(user, accounts, true)
			
			return response
		} else {
			return map[string]interface{}{"message": "not valid values"}
		}
}

func GetUser(id string, jwt string) map[string]interface{}  {
	isValid := helpers.ValidateToken(id, jwt)

	if isValid {
		db := helpers.ConnectDB()
		user := &interfaces.User{}

		dbRresult := db.Where("id = ?", id).First(&user)
		if errors.Is(dbRresult.Error, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found"}
		}

		// Find account for the user
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, false)
		return response
	} else {
		return map[string]interface{}{"message": "Not valid token"}
	}
}