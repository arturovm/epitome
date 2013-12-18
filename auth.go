package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	SessionToken string    `json:"session_token"`
	AppName      string    `json:"app_name"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetUserForSessionToken(token string) (*User, error, int) {
	var u User
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		return nil, errors.New("Couldn't connect to database"), 500
	}
	defer DB.Close()
	err = DB.QueryRow("select username from sessions where session_token = ?", token).Scan(&u.Username)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid session token"), 401
	}
	err = DB.QueryRow("select * from users where username = ?", u.Username).Scan(&u.Id, &u.Username, &u.DisplayName, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, errors.New("Couldn't retrieve logged in user"), 500
	}
	return &u, nil, 0
}

func PostSessions(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")
	appName := req.FormValue("appname")
	if username == "" || password == "" {
		WriteJSONError(w, http.StatusBadRequest, "Not enough parameters to log in")
		return
	}
	username = strings.ToLower(username)
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
		return
	}
	defer DB.Close()
	var u User
	err = DB.QueryRow("select * from users where username = ?", username).Scan(&u.Id, &u.Username, &u.DisplayName, &u.PasswordHash, &u.Role)
	if err == sql.ErrNoRows {
		WriteJSONError(w, http.StatusUnauthorized, "Incorrect username or password.")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Incorrect username or password.")
		return
	}
	unixTimestamp := time.Now().Unix()
	md5Hash := md5.New()
	io.WriteString(md5Hash, username+":"+strconv.Itoa(int(unixTimestamp)))
	sessionToken := hex.EncodeToString(md5Hash.Sum(nil))
	_, err = DB.Exec("insert into sessions values (null, ?, ?, ?, datetime(?, 'unixepoch'))", username, sessionToken, appName, unixTimestamp)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't write to database")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"session_token": "` + sessionToken + `"}`))
	return
}

func DeleteSessions(w http.ResponseWriter, req *http.Request) {
	sessionToken := req.Header.Get("x-session-token")
	if sessionToken == "" {
		WriteJSONError(w, http.StatusBadRequest, "Session token not provided")
		return
	}
	u, err, code := GetUserForSessionToken(sessionToken)
	if err != nil {
		WriteJSONError(w, code, err.Error())
		return
	}
	token := req.URL.Query().Get(":token")
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	defer DB.Close()
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
		return
	}
	_, err = DB.Exec("delete from sessions where session_token=? and username=?", token, u.Username)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Unable to delete session")
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
