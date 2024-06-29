package main

import (
	"fmt"
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

func verifyJWTToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
