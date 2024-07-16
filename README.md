# CapyBank

Bank simulation API built with Go's standard library. Thought of as simple  
API to demonstrate Go's capabilities in building a RESTful API.

## Features

- Create an account (sign up & login)
- Transactions (proper validations & custom error handling)
  - Deposit
  - Withdraw
  - Transfer

## Endpoints

- GET /transactions (latest transactions with limit & offset support)
- POST /signup
- POST /login
- POST /transactions (transaction handler for deposit, withdraw, and transfer)

## Technical Details

- API built with Go's 1.22 latest stdlib features (net/http routing)
- PostgresSQL for data storage (transactions, accounts)
  - Passwords are hashed with bcrypt
- Authentication with JWT & auth middleware
- Custom error handling with error codes & JSON responses
- Custom structured logging for both server & client errors and requests
- Environment variables and command-line flags for configuration
