// quick check how bcrypt salts the passwords
package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	storedPasswords := map[string][]byte{
		"massive-100-mean-amazed": []byte{}, // a good enough password
		"HelloKitty99!":           []byte{}, // a password many people use like
	}

	for k := range storedPasswords {
		pwd, err := bcrypt.GenerateFromPassword([]byte(k), 10)
		if err != nil {
			log.Fatal(err)
		} else {
			storedPasswords[k] = pwd
		}
	}

	for original, hashed := range storedPasswords {
		fmt.Println(original, "\n", string(hashed), len(hashed), "\n")
	}
}
