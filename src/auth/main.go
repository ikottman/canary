package main

import (
	"fmt"
	"log"

	"github.com/ikottman/canary/auth"
)

func main() {
	token, err := auth.CreateJwt()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Token:", token)

	validated := auth.ValidateJwt(token)
	fmt.Println("Valid:", validated)
}
