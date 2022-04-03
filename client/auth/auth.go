package auth

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
)

/**
* Creates an RSA 256 encoded JWT with no claims.
* You must specify the env var CANARY_PRIVATE_KEY
* with a base64 encoded RSA private key (to deal with newlines)
**/
func CreateJwt() (string, error) {
	privateKey, _ := b64.RawStdEncoding.DecodeString(os.Getenv("CANARY_PRIVATE_KEY"))
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	token, err := jwt.New(jwt.SigningMethodRS256).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}
