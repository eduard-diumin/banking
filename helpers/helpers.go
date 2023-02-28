package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/eduard-diumin/banking/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}

func ConnectDB() *gorm.DB {
	dsn := "host=localhost user=banking password=banking dbname=banking_development port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	HandleErr(err)
	return db
}

func Validation(values []interfaces.Validation) bool {
	username := regexp.MustCompile(`^([A-Za-z0-9]{5,})+$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z0-9]+$`)

	for i := 0; i < len(values); i++{
		switch values[i].Valid {
			case "username":
				if !username.MatchString(values[i].Value){
					return false
				}
			case "email":
				if !email.MatchString(values[i].Value){
					return false
				}
			case "password":
				if len(values[i].Value) < 5{
					return false
				}
		}
	}
	return true
}

func PanicHandle(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func ()  {
			error := recover()
			if error != nil {
				log.Println(error)
				resp := interfaces.ErrResponse{Message: "Internal server error"}
				json.NewEncoder(w).Encode(resp)
			}	
		}()
		next.ServeHTTP(w,r)
	})
}

func ValidateToken(id string, jwtToken string) bool  {
	clearJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(clearJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)
	var userID, _ = strconv.ParseFloat(id, 8)
	if token.Valid && tokenData["user_id"] == userID {
		return true
	} else {
		return false
	}
}
