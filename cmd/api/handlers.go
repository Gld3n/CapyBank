package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user *User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	fmt.Println(user.Username, user.Password)
}

func (app *application) depositHandler(w http.ResponseWriter, r *http.Request) {
	var req *Transaction

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.logger.Info(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if req.Type == Withdrawal {
		req.Amount -= req.Amount * 2
	}

	if err := app.transactions.makeNewTransaction(req); err != nil {
		app.serverError(w, r, err)
	}
}

//func (app *application) transferHandler(w http.ResponseWriter, r *http.Request) {
//}
//func (app *application) withdrawalHandler(w http.ResponseWriter, r *http.Request) {
//}
