package main

import "time"

type RequestCreateUser struct {
	Role     Role   `json:"role"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestLoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUser struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBUser struct {
	ID             int       `json:"id"`
	Role           Role      `json:"role"`
	Email          string    `json:"email"`
	Balance        float64   `json:"balance"`
	Fullname       string    `json:"fullname"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Role string

const (
	Admin  Role = "admin"
	Normal      = "normal"
)

type Operation string

const (
	Deposit    Operation = "deposit"
	Transfer             = "transfer"
	Withdrawal           = "withdrawal"
)

type Transaction struct {
	UserID       int       `json:"user_id"`
	Amount       float64   `json:"amount"`
	Type         Operation `json:"transaction_type"`
	TargetUserID int       `json:"target_user_id"`
}
