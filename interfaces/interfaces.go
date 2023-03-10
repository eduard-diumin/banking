package interfaces

import (
	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

type ResponseAccount struct {
	ID uint
	Name string
	Balance int
}

type ResponseUser struct {
	ID uint
	FirstName string
	LastName  string
	Email     string
	Account []ResponseAccount
}

type Validation struct {
	Value string
	Valid string
}

type ErrResponse struct {
	Message string
}

type Transactions struct {
	gorm.Model
	From uint
	To uint
	Amount int
}