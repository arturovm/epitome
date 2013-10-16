package main

import (
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + errorMessage + `"}`))
}
