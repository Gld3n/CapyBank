package main

type Role string

const (
	Admin  Role = "admin"
	Normal      = "normal"
)

type User struct {
	id       int
	role     Role
	email    string
	balance  float64
	fullname string
	username string
	password string
}

type Operation string

const (
	Deposit    Operation = "deposit"
	Transfer             = "transfer"
	Withdrawal           = "withdrawal"
)

type Transaction struct {
	UserID int       `json:"user_id"`
	Amount float64   `json:"amount"`
	Type   Operation `json:"transaction_type"`
}
