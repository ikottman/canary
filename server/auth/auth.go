package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
)

/**
* Validates an RSA 256 encoded JWT.
* You must specify the env var CANARY_PUBLIC_KEY
* with a base64 encoded RSA public key (to deal with newlines)
**/
func ValidateJwt(token string) bool {
	publicKey, _ := b64.RawStdEncoding.DecodeString(os.Getenv("CANARY_PUBLIC_KEY"))
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		fmt.Println("Failed to parse public key, error:", err)
		return false
	}

	parsed, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return false, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})

	if err != nil || parsed == nil {
		fmt.Println("Failed to validate token, error:", err)
		return false
	}

	return true
}
