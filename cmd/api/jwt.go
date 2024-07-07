package main

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func createJWTToken(user *DBUser) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"iss": "capybank",
		"aud": user.Role,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func verifyJWTToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return err
}
