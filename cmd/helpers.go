package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {

	app.logger.Err(err).Msg(fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()))

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

type RapidResponse struct {
	Get        string        `json:"get"`
	Parameters interface{}   `json:"parameters"`
	Errors     interface{}   `json:"errors"`
	Results    int           `json:"results"`
	Paging     interface{}   `json:"paging"`
	Response   []interface{} `json:"response"`
}

func (app *application) parseRapidError(resp RapidResponse) error {
	return errors.New(resp.Errors.(string))
}
