package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /login", http.HandlerFunc(app.loginHandler))
	mux.Handle("POST /signup", http.HandlerFunc(app.signupHandler))
	mux.Handle("POST /transaction/deposit", app.requireAuthentication(http.HandlerFunc(app.depositHandler)))
	//mux.Handle("POST /transaction/transfer", http.HandlerFunc(app.transferHandler))
	//mux.Handle("POST /transaction/withdrawal", http.HandlerFunc(app.withdrawalHandler))

	return app.logRequest(mux)
}
