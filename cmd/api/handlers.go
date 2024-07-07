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

	dbUser, err := app.users.getUserByUsername(reqLogin.Username)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			app.clientError(w, err, http.StatusUnauthorized)
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
		app.logger.Info(err.Error())
		app.clientError(w, err, http.StatusBadRequest)
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
		app.serverError(w, r, err)
	}
}
