package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	//"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
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
	if verboseMode == true {
		log.SetOutput(os.Stdout)
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		log.Print("Received request at '/api/users'\n" + string(reqB) + "\n\n\n")
		log.SetOutput(os.Stderr)
	}
	uAuth, err, _ := GetUserForSessionToken(req.Header.Get("x-session-token"))
	if *UserPreferences.NewUserPermissions != PublicRole && (uAuth == nil || uAuth.Role > *UserPreferences.NewUserPermissions) {
		WriteJSONError(w, http.StatusUnauthorized, "Insufficient permissions to create new accounts")
		return
	}
	var role UserRole
	role = NormalRole
	if req.FormValue("role") == "admin" {
		role = AdminRole
	}
	username := req.FormValue("username")
	password := req.FormValue("password")
	if username == "" || password == "" {
		WriteJSONError(w, http.StatusBadRequest, "Insufficient parameters")
		return
	}
	usernameLower := strings.ToLower(username)
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
}

func GetUser(w http.ResponseWriter, req *http.Request) {
}

func PutUsers(w http.ResponseWriter, req *http.Request) {
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
}
