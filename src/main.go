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

	validated := auth.ValidateJwt(token)
	fmt.Println("Valid:", validated)
}
