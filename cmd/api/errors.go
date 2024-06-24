package main

import (
	"log/slog"
	"net/http"
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
