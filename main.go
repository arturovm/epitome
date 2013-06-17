package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strings"
)

func createTables() {
	DB, _ := sql.Open("sqlite3", "db.db")
	DB.Exec("create table subscriptions (id integer primary key, url text, name text)")
	DB.Exec("create table articles (id integer primary key, url text, name text, author text, published integer, parent_id integer references subscriptions(id), body text, read bool)")
	DB.Exec("create table favorites (id integer primary key, url text, name text, author text, published integer, body text)")
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	port := flag.String("port", "8684", "The port in which the server will listen and serve")
	flag.StringVar(port, "p", "8684", "The port in which the server will listen and serve")
	flag.Parse()
	if _, err := os.Open("db.db"); os.IsNotExist(err) {
		if _, err := os.Create("db.db"); err != nil {
			log.Fatal(err)
		}
		createTables()
	}
	m := RegisterRoutes()
	http.Handle("/api/", m)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	if !strings.HasPrefix(*port, ":") {
		*port = ":" + *port
	}
	log.Print("Starting server on port " + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
