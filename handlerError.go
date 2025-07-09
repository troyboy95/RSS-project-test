package main

import (
	"net/http"
)

func handlerError(w http.ResponseWriter, r *http.Request) {
	// Log the error
	respondWithError(w, 400, "Sth went wrong, please try again later")
}
