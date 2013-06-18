package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var ExePath string

func createTables() {
	DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
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
	var err error
	ExePath, err = filepath.EvalSymlinks(os.Args[0])
	ExePath = filepath.Dir(ExePath)
	if err != nil {
		log.Fatal("Unable to resolve path to exectuable")
	}
	if _, err := os.Open(ExePath + "/db.db"); os.IsNotExist(err) {
		if _, err := os.Create(ExePath + "/db.db"); err != nil {
			log.Fatal(err)
		}
		createTables()
	}
	m := RegisterRoutes()
	http.Handle("/api/", m)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(ExePath+"/static"))))
	if !strings.HasPrefix(*port, ":") {
		*port = ":" + *port
	}
	log.Print("Starting server on port " + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
