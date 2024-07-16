package main

import (
	"database/sql"
	"errors"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type QueryRower interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func getUserByUsername(q QueryRower, username string) (DBUser, error) {
	stmt := `SELECT id, user_role, email, balance, fullname, username, password, created_at, updated_at FROM users WHERE username = $1`

	var dbu DBUser

	err := q.QueryRow(stmt, username).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.Balance, &dbu.Fullname, &dbu.Username, &dbu.HashedPassword, &dbu.CreatedAt, &dbu.UpdatedAt)
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

func (u *UserModel) updateUser(userId int, userAttr string, newValue interface{}) error {
	stmt := `UPDATE users SET $1 = $2 WHERE id = $3`

	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	if _, err = tx.Exec(stmt, userAttr, newValue, userId); err != nil {
		return err
	}

	return nil
}

func updateBalance(tx *sql.Tx, balance float64, userId int) error {
	_, err := tx.Exec(`UPDATE users SET balance = $1 WHERE id = $2`, balance, userId)
	return err
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) insertTransaction(tx *sql.Tx, tr *Transaction) error {
	stmt := `INSERT INTO transactions (user_id, amount, transaction_type, target_user_id, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(stmt, tr.UserID, tr.Amount, tr.Type, tr.TargetUserID, time.Now())
	return err
}

func (t *TransactionModel) getLatestTransactions(limit, offset int) ([]Transaction, error) {
	tx, err := t.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var transactions []Transaction

	rows, err := tx.Query(`SELECT user_id, amount, transaction_type, target_user_id FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var t Transaction

		if err = rows.Scan(&t.UserID, &t.Amount, &t.Type, &t.TargetUserID); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
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
		return ErrInvalidOperationType
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
	if tr.UserID == tr.TargetUserID {
		return ErrSameUserTransaction
	}

	var targetExists int
	err := tx.QueryRow(`SELECT 1 FROM users WHERE id = $1`, tr.TargetUserID).Scan(&targetExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	newBalance := balance - tr.Amount
	if newBalance < 0 {
		return ErrInsufficientFunds
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, tr.Amount, tr.TargetUserID)
	if err != nil {
		return err
	}

	return updateBalance(tx, newBalance, tr.UserID)
}

func withdraw(tx *sql.Tx, tr *Transaction, balance float64) error {
	newBalance := balance - tr.Amount
	if newBalance < 0 {
		return ErrInsufficientFunds
	}

	return updateBalance(tx, newBalance, tr.UserID)
}
