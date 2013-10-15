package main

import (
	//"database/sql"
	//"encoding/json"
	//_ "github.com/mattn/go-sqlite3"
	"net/http"
)

type UserRole int

const (
	AdminRole UserRole = iota
	NormalRole
	PublicRole
)

type User struct {
	Id           int      `json:"id"`
	Username     string   `json:"username"`
	PasswordHash string   `json:"-"`
	Role         UserRole `json:"role"`
}

func PostUser(w http.ResponseWriter, req *http.Request) {
}

func GetUser(w http.ResponseWriter, req *http.Request) {
}

func PutUsers(w http.ResponseWriter, req *http.Request) {
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
}
