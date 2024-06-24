package main

import (
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
