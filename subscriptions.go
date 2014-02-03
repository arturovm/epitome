package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Subscription struct {
	Id     int    `json:"id"`
	UserId int    `json:"-"`
	Url    string `json:"url"`
	Name   string `json:"name"`
}

func removeTrail(rawurl string) string {
	if rawurl[len(rawurl)-1] == '/' {
		return rawurl[:len(rawurl)-1]
	}
	return rawurl
}

func findRSSURL(rawurl string) (string, error) {
	res, err := http.Get(rawurl)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	doc, err := gokogiri.ParseXml(body)
	defer doc.Free()
	if (doc.Root().Name() != "rss" && doc.Root().Name() != "feed") || err != nil {
		doc, _ := gokogiri.ParseHtml(body)
		defer doc.Free()
		doc.RecursivelyRemoveNamespaces()
		nodes, err := doc.Search("//link")
		if err != nil {
			return "", err
		}
		for _, v := range nodes {
			if v.Attribute("rel").Value() == "alternate" || v.Attribute("rel").Value() == "feed" {
				u, _ := url.Parse(v.Attribute("href").Value())
				if u.IsAbs() {
					return u.String(), nil
				} else {
					baseU, _ := url.Parse(rawurl)
					u.Scheme = baseU.Scheme
					u.Host = baseU.Host
					return u.String(), nil
				}
			}
		}
	} else {
		return res.Request.URL.String(), nil
	}
	return "", errors.New("Feed URL not found")
}

func findRSSTitle(rssUrl string) (string, error) {
	res, err := http.Get(rssUrl)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	doc, err := gokogiri.ParseXml(body)
	defer doc.Free()
	doc.RecursivelyRemoveNamespaces()
	nodes, err := doc.Search("//title")
	if err != nil {
		return "", err
	}
	if len(nodes) > 0 {
		return nodes[0].Content(), nil
	}
	return "", nil
}

func PostSubscription(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.SetOutput(os.Stdout)
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		log.Print("Received request at '/api/subscriptions'\n" + string(reqB) + "\n\n\n")
		log.SetOutput(os.Stderr)
	}
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
	if rawurl := req.FormValue("url"); rawurl != "" {
		rssUrl, err := findRSSURL(rawurl)
		if err != nil {
			WriteJSONError(w, http.StatusNotFound, "Could not find feed.")
			log.Printf("POST Subscription error (Find URL error): %s", err.Error())
			return
		} else {
			title, err := findRSSTitle(rssUrl)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(""))
				log.Printf("POST Subscription error (Find title error): %s", err.Error())
				return
			}
			DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
			defer DB.Close()
			var sub Subscription
			rowErr := DB.QueryRow("select * from subscriptions where url=?", rssUrl).Scan(&sub.Id, &sub.Url, &sub.Name)
			if rowErr == nil {
				rowErr = DB.QueryRow("select * from user_subscriptions where subscription_id=? and user_id=?", sub.Id, u.Id).Scan()
				if rowErr == sql.ErrNoRows {
					DB.Exec("insert into user_subscriptions values (null, ?, ?)", sub.Id, u.Id)
					w.WriteHeader(http.StatusCreated)
					return
				} else {
					WriteJSONError(w, http.StatusConflict, "Feed already exists")
					return
				}
			} else if rowErr == sql.ErrNoRows {
				DB.Exec("insert into subscriptions values (null, ?, ?)", rssUrl, title)
				DB.QueryRow("select id from subscriptions where url=?", rssUrl).Scan(&sub.Id)
				DB.Exec("insert into user_subscriptions values (null, ?, ?)", sub.Id, u.Id)
				w.WriteHeader(http.StatusCreated)
				return
			}
		}
	} else {
		WriteJSONError(w, http.StatusBadRequest, "Insufficient parameters: URL was not provided")
		return
	}
}

func GetSubscriptions(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.SetOutput(os.Stdout)
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		log.Print("Received request at '/api/subscriptions'\n" + string(reqB) + "\n\n\n")
		log.SetOutput(os.Stderr)
	}
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
	global := req.FormValue("global")
	DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
	defer DB.Close()

	subs := make([]Subscription, 0)

	switch global {
	case "true":
		if u.Role != AdminRole {
			WriteJSONError(w, http.StatusUnauthorized, "You don't have enough permissions to view global subscriptions")
			return
		}
		rows, err := DB.Query("select * from subscriptions")
		if err != nil {
			log.Printf("GET Subscriptions error (Database error): %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var sub Subscription
			rows.Scan(&sub.Id, &sub.Url, &sub.Name)
			subs = append(subs, sub)
		}
		rows.Close()
	default:
		rows, err := DB.Query("select subscription_id from user_subscriptions where user_id=?", u.Id)
		if err != nil {
			log.Printf("GET Subscriptions error (Database error): %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var sub Subscription
			rows.Scan(&sub.Id)
			DB.QueryRow("select url, name from subscriptions where id=?", sub.Id).Scan(&sub.Url, &sub.Name)
			subs = append(subs, sub)
		}
		rows.Close()
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	encErr := enc.Encode(subs)
	if encErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		log.Printf("GET Subscriptions error (JSON encoding error): %s", err.Error())
	}
}

func DeleteSubscription(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.SetOutput(os.Stdout)
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		log.Print("Received request at '/api/subscriptions/:id'\n" + string(reqB) + "\n\n\n")
		log.SetOutput(os.Stderr)
	}
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
	id := req.URL.Query().Get(":id")
	global := req.FormValue("global")
	DB, _ := sql.Open("sqlite3", ExePath+"/db.db")
	defer DB.Close()
	var dummy int
	err = DB.QueryRow("select id from subscriptions where id = ?", id).Scan(&dummy)
	switch {
	case err == sql.ErrNoRows:
		WriteJSONError(w, http.StatusNotFound, "Subscription does not exist")
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		log.Printf("DELETE subscription error (Database error): %s", err.Error())
	default:
		if global == "true" {
			if u.Role != AdminRole {
				WriteJSONError(w, http.StatusUnauthorized, "You don't have enough permissions to delete global subscriptions")
				return
			}
			DB.Exec("delete from subscriptions where id=?", id)
			DB.Exec("delete from user_subscriptions where subscription_id=?", id)
			w.WriteHeader(http.StatusOK)
			return
		}
		var sub Subscription
		err := DB.QueryRow("select id from user_subscriptions where subscription_id=? and user_id=?", id, u.Id).Scan(&sub.Id)
		if err == sql.ErrNoRows {
			WriteJSONError(w, http.StatusNotFound, "Subscription doesn't exist")
			return
		}
		DB.Exec("delete from user_subscriptions where id=?", sub.Id)
		w.WriteHeader(http.StatusOK)
		return
	}
}
