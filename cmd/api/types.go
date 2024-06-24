package main

type Role string

const (
	Admin  Role = "admin"
	Normal      = "normal"
)

type User struct {
	Id       int     `json:"id"`
	Role     Role    `json:"role"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	Fullname string  `json:"fullname"`
	Username string  `json:"username"`
	Password string  `json:"password"`
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
