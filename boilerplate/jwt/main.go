package main

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("you'llneverfindout")

func main() {

	// Create
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user"] = "38"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatal(err)
	}

	// verify
	token, err = jwt.Parse(tokenString, signingKey)
	if err != nil {
		fmt.Println("Unable to verify:", err)
	}
	fmt.Println(token.Valid)
	fmt.Println(token.Raw)
}

// callback for jwt.Parse .. callback.. callback!.. callback?!
func signingKey(token *jwt.Token) (interface{}, error) {

	// only process HMAC signing
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	fmt.Println(token.Claims["user"])
	fmt.Println(time.Unix(int64(token.Claims["exp"].(float64)), 0)) //hopefully this can be nicer in a future version
	// TODO: add error checking to the type assertions above.
	return mySigningKey, nil
}
