package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("new request", slog.String("ip", ip), slog.String("protocol", proto), slog.String("method", method), slog.String("uri", uri))

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")[len("Bearer "):]
		if headerToken == "" {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		_, err := verifyJWTToken(headerToken)
		// TODO: create custom error json response struct and handle custom messages
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				app.clientError(w, http.StatusUnauthorized)
				return
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				app.clientError(w, http.StatusBadRequest)
			default:
				app.serverError(w, r, err)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
