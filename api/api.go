package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/eduard-diumin/banking/helpers"
	"github.com/eduard-diumin/banking/interfaces"
	"github.com/eduard-diumin/banking/useraccounts"
	"github.com/eduard-diumin/banking/users"
	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email string
	Password string
}

type TransactionBody struct {
	UserID uint
	From uint
	To uint
	Amount int
}

func readBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	helpers.HandleErr(err)

	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter)  {
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		reps := interfaces.ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(reps)
	}	
}

func login(w http.ResponseWriter, r *http.Request) {
	// Redy body
	body := readBody(r)

	// Handle Login
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	login := users.Login(formattedBody.Username, formattedBody.Password)

	//Prepare response
	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	// Redy body
	body := readBody(r)

	// Handle Register
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)

	//Prepare response
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	userID := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userID, auth)
	apiResponse(user, w)
}

func transaction(w http.ResponseWriter, r *http.Request)  {
	// Redy body
	body := readBody(r)
	auth := r.Header.Get("Authorization")

	// Handle Register
	var formattedBody TransactionBody
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	transaction := useraccounts.Transaction(formattedBody.UserID, formattedBody.From, formattedBody.To, formattedBody.Amount, auth)

	//Prepare response
	apiResponse(transaction, w)
}

func StartApi()  {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandle)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	fmt.Println("App is working on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}