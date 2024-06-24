package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) depositHandler(w http.ResponseWriter, r *http.Request) {
	var req *Transaction

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.logger.Info(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
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
