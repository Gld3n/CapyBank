package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) getUserByUsername(username string) (DBUser, error) {
	stmt := `SELECT username, password FROM users WHERE username = $1`

	var dbu DBUser

	err := u.DB.QueryRow(stmt, username).Scan(&dbu.Username, &dbu.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DBUser{}, ErrNoRecord
		}
		return DBUser{}, err
	}

	return dbu, nil
}

func (u *UserModel) createUser(usr *RequestCreateUser) error {
	stmt := `INSERT INTO users (fullname, username, password, email, user_role) VALUES ($1, $2, $3, $4, $5)`

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	if _, err = tx.Exec(stmt, usr.Fullname, usr.Username, usr.Password, usr.Email, usr.Role); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) processNewTransaction(tr *Transaction) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var userBalance float64
	_ = tx.QueryRow(`SELECT balance FROM users WHERE id = $1`, tr.UserID).Scan(&userBalance)

	switch tr.Type {
	case Deposit:
		if err = deposit(tx, tr, userBalance); err != nil {
			return err
		}
	case Transfer:
	case Withdrawal:
	default:
		return fmt.Errorf("operation type %s not valid", tr.Type)
	}

	_, err = tx.Exec(`INSERT INTO transactions (user_id, amount, transaction_type, created_at) VALUES ($1, $2, $3, $4)`, tr.UserID, tr.Amount, tr.Type, time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func deposit(tx *sql.Tx, tr *Transaction, balance float64) error {
	newBalance := balance + tr.Amount

	_, err := tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2`, newBalance, tr.UserID)
	return err
}
