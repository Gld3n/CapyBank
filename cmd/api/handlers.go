package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reqLogin *RequestLoginUser
	if err := json.NewDecoder(r.Body).Decode(&reqLogin); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	dbUser, err := getUserByUsername(app.users.DB, reqLogin.Username)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			app.clientError(w, ErrUserNotFound, http.StatusNotFound)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	if !isPasswordEqualToHash(dbUser.HashedPassword, []byte(reqLogin.Password)) {
		app.clientError(w, err, http.StatusUnauthorized)
		return
	}

	token, err := createJWTToken(&dbUser)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	payload := make(map[string]string)
	payload["token"] = token

	if err = json.NewEncoder(w).Encode(payload); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var reqUser *RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	hash, err := hashPassword([]byte(reqUser.Password))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	reqUser.Password = hash

	if err = app.users.createUser(reqUser); err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) transactionHandler(w http.ResponseWriter, r *http.Request) {
	var reqTr *Transaction

	if err := json.NewDecoder(r.Body).Decode(&reqTr); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if reqTr.Amount <= 0 {
		app.clientError(w, ErrInvalidAmount, http.StatusBadRequest)
		return
	}

	if err := app.transactions.processNewTransaction(reqTr); err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			app.clientError(w, err, http.StatusNotFound)
		case errors.Is(err, ErrInsufficientFunds):
			app.clientError(w, err, http.StatusPaymentRequired)
		case isTransactionBadRequest(err):
			app.clientError(w, err, http.StatusBadRequest)
		default:
			app.serverError(w, r, err)
		}
	}
}

func (app *application) listTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	limit := 10
	offset := 0

	queryValues := r.URL.Query()
	limitStr := queryValues.Get("limit")
	offsetStr := queryValues.Get("offset")

	if limitStr != "" {
		err := validateLimit(limitStr, &limit)
		if err != nil {
			app.clientError(w, err, http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		err := validateOffset(offsetStr, &offset)
		if err != nil {
			app.clientError(w, err, http.StatusBadRequest)
			return
		}
	}

	transactions, err := app.transactions.getLatestTransactions(limit, offset)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	payload := make(map[string][]Transaction)
	payload["latest_transactions"] = transactions

	if err = json.NewEncoder(w).Encode(payload); err != nil {
		app.serverError(w, r, err)
	}
}
