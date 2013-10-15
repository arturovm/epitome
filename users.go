package main

import (
	//"database/sql"
	//"encoding/json"
	//_ "github.com/mattn/go-sqlite3"
	//"net/http"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}
