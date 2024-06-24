package main

import (
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) makeNewTransaction(tr *Transaction) error {
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
