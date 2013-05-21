package main

import (
	"net/http"
	"log"
)

func main() {
	m := RegisterRoutes()
	http.Handle("/api/",  m)
	http.Handle("/", http.FileServer(http.Dir("")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
