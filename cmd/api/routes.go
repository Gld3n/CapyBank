package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /transactions", app.requireAuthentication(http.HandlerFunc(app.listTransactionsHandler)))

	mux.Handle("POST /login", http.HandlerFunc(app.loginHandler))
	mux.Handle("POST /signup", http.HandlerFunc(app.signupHandler))
	mux.Handle("POST /transaction", app.requireAuthentication(http.HandlerFunc(app.transactionHandler)))

	return app.logRequest(mux)
}
