package main

import (
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds for this operation")
	ErrInvalidAmount     = errors.New("amount must be a positive number")
	ErrNoRecord          = errors.New("no matching record found")
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(
		err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
	)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, err error, status int) {
	app.logger.Error(err.Error())
	http.Error(w, http.StatusText(status), status)
}
