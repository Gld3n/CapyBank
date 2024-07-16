package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrInsufficientFunds    = errors.New("insufficient funds for this operation")
	ErrInvalidAmount        = errors.New("amount must be a positive number")
	ErrInvalidLimit         = errors.New("invalid limit value")
	ErrInvalidOffset        = errors.New("invalid offset value")
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrLimitExceeded        = errors.New("maximum allowed limit is 100")
	ErrNoRecord             = errors.New("no matching record found")
	ErrNoTargetSpecified    = errors.New("no target specified for transfer transaction")
	ErrSameUserTransaction  = errors.New("same user transfer not available")
	ErrUserNotFound         = errors.New("target user requested not found")
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
	response := map[string]string{"error": err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
