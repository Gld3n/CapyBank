package main

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func isPasswordEqualToHash(hash string, pwd []byte) bool {
	byteHash := []byte(hash)
	if err := bcrypt.CompareHashAndPassword(byteHash, pwd); err != nil {
		return false
	}

	return true
}
