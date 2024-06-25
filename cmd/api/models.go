package main

import (
	"database/sql"
	"errors"
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

func (u *UserModel) createNewUser(usr *RequestCreateUser) error {
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

func (t *TransactionModel) createNewTransaction(tr *Transaction) error {
	tx, err := t.DB.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	_, err = tx.Exec(`INSERT INTO transactions (user_id, amount, transaction_type, created_at) VALUES ($1, $2, $3, $4)`, tr.UserID, tr.Amount, tr.Type, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2;`, tr.Amount, tr.UserID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
