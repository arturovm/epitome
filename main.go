package main

import (
	"net/http"
	"log"
)

func main() {
	m := RegisterRoutes()
	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
