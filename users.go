package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	//"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strings"
)

type UserRole int

const (
	AdminRole UserRole = iota
	NormalRole
	PublicRole
)

type User struct {
	Id           int      `json:"id"`
	Username     string   `json:"-"`
	DisplayName  string   `json:"username"`
	PasswordHash string   `json:"-"`
	Role         UserRole `json:"role"`
}

func PostUser(w http.ResponseWriter, req *http.Request) {
	//TODO: check current user here and user role. For now, only check if open to public
	username := req.FormValue("username")
	usernameLower := strings.ToLower(username)
	password := req.FormValue("password")
	var role UserRole
	switch req.FormValue("role") {
	case "admin":
		role = AdminRole
	case "normal":
		role = NormalRole
	}
	if username == "" || password == "" {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Insufficient parameters"}`))
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Could not hash password"}`))
		return
	}
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Could not connect to database"}`))
		return
	}
	defer DB.Close()

	DB.Exec("insert into users values (null, ?, ?, ?, ?)", usernameLower, username, hash, role)
	w.WriteHeader(http.StatusCreated)
}

func GetUser(w http.ResponseWriter, req *http.Request) {
}

func PutUsers(w http.ResponseWriter, req *http.Request) {
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
}
