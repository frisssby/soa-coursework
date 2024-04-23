package jwt

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtPublic *rsa.PublicKey
var jwtPrivate *rsa.PrivateKey

func LoadKeys(jwtprivateFile *string, jwtPublicFile *string) error {
	if private, err := os.ReadFile(*jwtprivateFile); err != nil {
		return err
	} else if jwtPrivate, err = jwt.ParseRSAPrivateKeyFromPEM(private); err != nil {
		return err
	}
	if public, err := os.ReadFile(*jwtPublicFile); err != nil {
		return err
	} else if jwtPublic, err = jwt.ParseRSAPublicKeyFromPEM(public); err != nil {
		return err
	}
	return nil
}

func GenerateJWT(username string, expirationTime time.Time) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{"username": username},
	)
	tokenString, err := token.SignedString(jwtPrivate)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (username string, err error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtPublic, nil
	})
	if err != nil || !token.Valid {
		log.Println("parse token: ", err.Error())
		return "", err
	}
	username = claims["username"].(string)
	return username, nil
}
