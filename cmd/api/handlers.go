package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var reqLogin *RequestLoginUser
	if err := json.NewDecoder(r.Body).Decode(&reqLogin); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dbUser, err := app.users.getUserByUsername(reqLogin.Username)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			app.clientError(w, http.StatusUnauthorized)
			return
		} else {
			app.serverError(w, r, err)
			return
		}
	}

	if !isPasswordEqualToHash(dbUser.HashedPassword, []byte(reqLogin.Password)) {
		app.clientError(w, http.StatusUnauthorized)
		return
	}
	app.logger.Info("user authenticated successfully")
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var reqUser *RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		app.logger.Info(err.Error())
		app.clientError(w, http.StatusBadRequest)
	}

	hash, err := hashPassword([]byte(reqUser.Password))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	reqUser.Password = hash

	if err = app.users.createNewUser(reqUser); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("user created successfully")
}

func (app *application) depositHandler(w http.ResponseWriter, r *http.Request) {
	var reqTr *Transaction

	if err := json.NewDecoder(r.Body).Decode(&reqTr); err != nil {
		app.logger.Info(err.Error())
		app.clientError(w, http.StatusBadRequest)
	}

	if reqTr.Type == Withdrawal {
		reqTr.Amount -= reqTr.Amount * 2
	}

	if err := app.transactions.createNewTransaction(reqTr); err != nil {
		app.serverError(w, r, err)
	}
}

//func (app *application) transferHandler(w http.ResponseWriter, r *http.Request) {
//}
//func (app *application) withdrawalHandler(w http.ResponseWriter, r *http.Request) {
//}
