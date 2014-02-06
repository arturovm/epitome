package main

import (
	"net/http"
	"strings"
)

func WriteJSONError(w http.ResponseWriter, status int, errorMessage string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + errorMessage + `"}`))
}

func TestContentType(rawHeader *string, target string) bool {
	parts := strings.Split(*rawHeader, ";")
	if len(parts) > 0 {
		return strings.Trim(parts[0], " \t") == target
	}
	return false
}