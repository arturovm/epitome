package main

import (
	"net/http"
	"log"
)

func main() {
	m := RegisterRoutes()
	http.Handle("/api/",  m)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
