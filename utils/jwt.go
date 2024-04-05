package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var secretToken = os.Getenv("SECRET_TOKEN")

func GenerateJwtToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	webToken, err := token.SignedString([]byte(secretToken))

	if err != nil {
		return "", err
	}

	return webToken, nil
}

func VerifyTokenJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, isValid := token.Method.(*jwt.SigningMethodHMAC)
		if !isValid {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretToken), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := VerifyTokenJwt(tokenString)
	if err != nil {
		return nil, err
	}

	claims, isOk := token.Claims.(jwt.MapClaims)

	if isOk && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("INVALID TOKEN")
}
