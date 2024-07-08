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

type QueryRower interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func getUserByUsername(q QueryRower, username string) (DBUser, error) {
	stmt := `SELECT id, user_role, username, password FROM users WHERE username = $1`

	var dbu DBUser

	err := q.QueryRow(stmt, username).Scan(&dbu.ID, &dbu.Role, &dbu.Username, &dbu.HashedPassword)
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

	return tx.Commit()
}

func updateBalance(tx *sql.Tx, balance float64, userId int) error {
	_, err := tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2`, balance, userId)
	return err
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) insertTransaction(tx *sql.Tx, tr *Transaction) error {
	_, err := tx.Exec(`INSERT INTO transactions (user_id, amount, transaction_type, created_at) VALUES ($1, $2, $3, $4)`, tr.UserID, tr.Amount, tr.Type, time.Now())
	return err
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
		if err = transfer(tx, tr, userBalance); err != nil {
			return err
		}
	case Withdrawal:
		if err = withdraw(tx, tr, userBalance); err != nil {
			return err
		}
	default:
		return fmt.Errorf("operation type %s not valid", tr.Type)
	}

	if err = t.insertTransaction(tx, tr); err != nil {
		return err
	}

	return tx.Commit()
}

func deposit(tx *sql.Tx, tr *Transaction, balance float64) error {
	newBalance := balance + tr.Amount
	return updateBalance(tx, newBalance, tr.UserID)
}

func transfer(tx *sql.Tx, tr *Transaction, balance float64) error {
	/* TODO: validate...
	- Target user exists
	- Target user is not same user
	- JSON nil target user to avoid nil pointer dereference
	*/
	newBalance := balance - tr.Amount
	if newBalance < 0 {
		return ErrInsufficientFunds
	}

	targetUser, err := getUserByUsername(tx, *tr.TargetUserUsername)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, tr.Amount, targetUser.ID)

	return updateBalance(tx, newBalance, tr.UserID)
}

func withdraw(tx *sql.Tx, tr *Transaction, balance float64) error {
	newBalance := balance - tr.Amount
	if newBalance < 0 {
		return ErrInsufficientFunds
	}

	return updateBalance(tx, newBalance, tr.UserID)
}
