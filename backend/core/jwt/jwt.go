package jwt

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func Encode(claims jwt.Claims, key *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

// Decode expects claims to be a pointer.
func Decode(claims jwt.Claims, tokenString string, publicKey crypto.PublicKey) error {
	withSigningMethodRS256 := jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()})
	parser := jwt.NewParser(withSigningMethodRS256)

	_, err := parser.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return err
	}

	// For multiple services or identity providers, the audience and issuer should be verified, but in our case,
	// we only have one identity provider (the UI auth/file server) and one audience: 3rivers services.
	// We also maintain exclusive access to the private key.

	return nil
}

func GetJWT(authorizationHeader string) (string, error) {
	split := strings.Split(authorizationHeader, " ")
	if len(split) != 2 || split[0] != "Bearer" || split[1] == "" {
		return "", fmt.Errorf("Authorization header is invalid >%s<", authorizationHeader)
	}

	return split[1], nil
}
