package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

func main() {
	if _, err := os.Open("db.db"); os.IsNotExist(err) {
		if _, err := os.Create("db.db"); err != nil {
			log.Fatal(err)
		}
		DB, _ := sql.Open("sqlite3", "db.db")
		DB.Exec("create table subscriptions (id int auto_increment pimary key, url string, name string)")
		if err := DB.Close(); err != nil {
			log.Fatal(err)
		}
	}
	m := RegisterRoutes()
	http.Handle("/api/", m)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
