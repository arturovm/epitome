package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
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
	}
	defer DB.Close()
	var u User
	err = DB.QueryRow("select * from users where username = ?", username).Scan(&u.Id, &u.Username, &u.DisplayName, &u.PasswordHash, &u.Role)
	if err == sql.ErrNoRows {
		WriteJSONError(w, http.StatusNotFound, "That user doesn't exist. Check your username.")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, "Wrong password")
		return
	}
	unixTimestamp := time.Now().Unix()
	md5Hash := md5.New()
	io.WriteString(md5Hash, username+":"+strconv.Itoa(int(unixTimestamp)))
	sessionToken := hex.EncodeToString(md5Hash.Sum(nil))
	_, err = DB.Exec("insert into sessions values (null, ?, ?, ?, datetime(?, 'unixexpoch'))", username, sessionToken, appName, unixTimestamp)
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
}
