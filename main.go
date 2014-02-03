package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ExePath         string
	CRON            *cron.Cron
	UserPreferences *Preferences
	verboseMode     bool
	verboseModeBody bool
)

const (
	MAX_ARTICLE_PROCS int = 4
)

func createTables() {
	DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
	DB.Exec("create table users (id integer primary key, username text, display_name text, password_hash text, role integer)")
	DB.Exec("create table sessions (id integer primary key, username text, session_token text, app_name text, created_at string)")
	DB.Exec("create table subscriptions (id integer primary key, url text, name text)")
	DB.Exec("create table user_subscriptions (id integer primary key, subscription_id integer references subscriptions(id), user_id integer references users(id))")
	DB.Exec("create table articles (id integer primary key, subscription_id integer references subscriptions(id), url text, title text, author text, published integer, body text, body_type text, summary text, summary_type text, read bool)")
	DB.Exec("create table user_read_articles (id integer primary key, user_id integer references users(id), article_id integer references articles(id), read bool)")
	DB.Exec("create table favorites (id integer primary key, url text, name text, author text, published integer, body text)")
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	port := flag.String("port", "8684", "The port in which the server will listen and serve.")
	flag.StringVar(port, "p", "8684", "The port in which the server will listen and serve.")
	flag.BoolVar(&verboseMode, "verbose", false, "Whether to log requests coming in to the Pond server and other details.")
	flag.BoolVar(&verboseModeBody, "log-body", false, "Whether to log the request body when verbose mode is on. If verbose mode is off, this option is a no-op. WARNING: If set to true while in verbose mode, memory usage could increase significantly.")
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
	ReloadPreferences()
	m := RegisterRoutes()
	http.Handle("/api/", m)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(ExePath+"/static"))))
	if !strings.HasPrefix(*port, ":") {
		*port = ":" + *port
	}
	log.Print("Starting server on port " + *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
