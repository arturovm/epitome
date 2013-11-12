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
	if UserPreferences.NewUserPermissions == PublicRole {
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
			WriteJSONError(w, http.StatusBadRequest, "Insufficient parameters")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, "Bcrypt Error")
			return
		}
		DB, err := sql.Open("sqlite3", ExePath+"/db.db")
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
			return
		}
		defer DB.Close()

		var u User
		err = DB.QueryRow("select * from users where username = ?", usernameLower).Scan(&u.Id, &u.Username, &u.DisplayName, &u.PasswordHash, &u.Role)
		if err == sql.ErrNoRows {
			_, err = DB.Exec("insert into users values (null, ?, ?, ?, ?)", usernameLower, username, hash, role)
			if err != nil {
				WriteJSONError(w, http.StatusInternalServerError, "Unable to write to database")
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			WriteJSONError(w, http.StatusConflict, "User already exists")
			return
		}
	} else {
		WriteJSONError(w, http.StatusUnauthorized, "You have insufficient permissions to create new accounts")
		return
	}
}

func GetUser(w http.ResponseWriter, req *http.Request) {
}

func PutUsers(w http.ResponseWriter, req *http.Request) {
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
}
