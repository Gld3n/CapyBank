package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
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

func validateLimit(queryLimit string, targetLimit *int) error {
	const minLimit = 1
	const maxLimit = 100

	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		return ErrInvalidLimit
	}
	if limit < minLimit || limit > maxLimit {
		return ErrLimitExceeded
	}

	*targetLimit = limit
	return nil
}

func validateOffset(queryOffset string, targetOffset *int) error {
	ofs, err := strconv.Atoi(queryOffset)
	if err != nil {
		return ErrInvalidOffset
	}

	*targetOffset = ofs
	return nil
}

func isTransactionBadRequest(err error) bool {
	return errors.Is(err, ErrInvalidOperationType) || errors.Is(err, ErrNoTargetSpecified) || errors.Is(err, ErrSameUserTransaction)
}
