package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

type ArticleContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Article struct {
	Id             int            `json:"id"`
	SubscriptionId int            `json:"subscription_id"`
	Url            string         `json:"url"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Published      time.Time      `json:"published_at"`
	Body           ArticleContent `json:"body"`
	Summary        ArticleContent `json:"summary"`
	Read           bool           `json:"read"`
}

func processQueryOptions(req *http.Request) *map[string]string {
	opts := make(map[string]string)
	opts["status"] = "all"
	if req.FormValue("status") == "read" {
		opts["status"] = "read"
	} else if req.FormValue("status") == "unread" {
		opts["status"] = "unread"
	}
	opts["order"] = "desc"
	if req.FormValue("order") == "asc" {
		opts["order"] = "asc"
	}
	opts["limit"] = "100"
	if req.FormValue("limit") != "" {
		opts["limit"] = req.FormValue("limit")
	}
	opts["since_id"] = ""
	if req.FormValue("since_id") != "" {
		opts["since_id"] = "and (datetime(published) >= (select datetime(published) from articles where id = " + req.FormValue("since_id") + ") and articles.id != " + req.FormValue("since_id") + ")"
	}
	opts["before_id"] = ""
	if req.FormValue("before_id") != "" {
		opts["before_id"] = "and (datetime(published) <= (select datetime(published) from articles where id = " + req.FormValue("before_id") + ") and articles.id != " + req.FormValue("before_id") + ")"
	}
	return &opts
}

func queryStringForRequest(req *http.Request, u *User, subsIds *[]string) string {
	opts := processQueryOptions(req)
	var query string
	switch (*opts)["status"] {
	case "read":
		query = fmt.Sprintf("select * from (select articles.id, subscription_id, url, title, author, datetime(published) as published, body, body_type, summary, summary_type, user_read_articles.read from articles inner join user_read_articles on articles.id = user_read_articles.article_id where subscription_id in (%s) and user_read_articles.user_id = %d %s %s order by datetime(published) desc limit %s) order by datetime(published) %s", strings.Join(*subsIds, ", "), u.Id, (*opts)["since_id"], (*opts)["before_id"], (*opts)["limit"], (*opts)["order"])
	case "unread":
		query = fmt.Sprintf("select * from (select articles.id, subscription_id, url, title, author, datetime(published) as published, body, body_type, summary, summary_type, articles.read from articles left outer join user_read_articles on articles.id = user_read_articles.article_id where subscription_id in (%s) and (user_read_articles.user_id = %d or user_read_articles.user_id is null) and user_read_articles.id is null %s %s order by datetime(published) desc limit %s) order by datetime(published) %s", strings.Join(*subsIds, ", "), u.Id, (*opts)["since_id"], (*opts)["before_id"], (*opts)["limit"], (*opts)["order"])
	case "all":
		query = fmt.Sprintf("select * from (select articles.id, subscription_id, url, title, author, datetime(published) as published, body, body_type, summary, summary_type, coalesce(articles.read, user_read_articles.read) from articles left outer join user_read_articles on articles.id = user_read_articles.article_id where subscription_id in (%s) and (user_read_articles.user_id = %d or user_read_articles.user_id is null) %s %s order by datetime(published) desc limit %s) order by datetime(published) %s", strings.Join(*subsIds, ", "), u.Id, (*opts)["since_id"], (*opts)["before_id"], (*opts)["limit"], (*opts)["order"])
	}
	return query
} // I AM TEH MASTAR OF SQLLLL

func GetAllArticles(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.Print("Received request at '/api/subscriptions/articles'")
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		io.WriteString(os.Stdout, string(reqB))
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
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
		return
	}
	var subs []Subscription
	subsRows, _ := DB.Query("select subscription_id from user_subscriptions where user_id = ?", u.Id)
	for subsRows.Next() {
		var sub Subscription
		subsRows.Scan(&sub.Id)
		subs = append(subs, sub)
	}
	subsRows.Close()
	var subsIds []string
	for _, v := range subs {
		subsIds = append(subsIds, fmt.Sprintln(v.Id))
	}
	var articles []Article
	rows, err := DB.Query(queryStringForRequest(req, u, &subsIds))
	for rows.Next() {
		var article Article
		var dateString string
		rows.Scan(&article.Id, &article.SubscriptionId, &article.Url, &article.Title, &article.Author, &dateString, &article.Body.Content, &article.Body.Type, &article.Summary.Content, &article.Summary.Type, &article.Read)
		article.Published, err = time.Parse("2006-01-02 15:04:05", dateString)
		article.Body.Content = html.EscapeString(article.Body.Content)
		article.Summary.Content = html.EscapeString(article.Summary.Content)
		articles = append(articles, article)
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	enc.Encode(articles)
}

func GetArticles(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.Print("Received request at '/api/subscriptions/:id/articles'")
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		io.WriteString(os.Stdout, string(reqB))
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
	subId := req.URL.Query().Get(":id")
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
		return
	}
	var id int
	err = DB.QueryRow("select id from user_subscriptions where user_id = ? and subscription_id = ?", u.Id, subId).Scan(&id)
	if err == sql.ErrNoRows {
		WriteJSONError(w, http.StatusNotFound, "Subscription doesn't exist")
		return
	}
	subsIds := []string{subId}
	var articles []Article
	rows, _ := DB.Query(queryStringForRequest(req, u, &subsIds))
	for rows.Next() {
		var article Article
		var dateString string
		rows.Scan(&article.Id, &article.SubscriptionId, &article.Url, &article.Title, &article.Author, &dateString, &article.Body.Content, &article.Body.Type, &article.Summary.Content, &article.Summary.Type, &article.Read)
		article.Published, err = time.Parse("2006-01-02 15:04:05", dateString)
		article.Body.Content = html.EscapeString(article.Body.Content)
		article.Summary.Content = html.EscapeString(article.Summary.Content)
		articles = append(articles, article)
	}
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")
	enc.Encode(articles)
}

func PutArticle(w http.ResponseWriter, req *http.Request) {
	if verboseMode == true {
		log.Print("Received request at '/api/subscriptions/:subid/articles/:artid'")
		reqB, _ := httputil.DumpRequest(req, verboseModeBody)
		io.WriteString(os.Stdout, string(reqB))
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
	subId := req.URL.Query().Get(":subid")
	artId := req.URL.Query().Get(":artid")
	status := req.FormValue("status")
	DB, err := sql.Open("sqlite3", ExePath+"/db.db")
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't connect to database")
		return
	}
	//var id int
	err = DB.QueryRow("select subscription_id from user_subscriptions where user_id = ? and subscription_id = ?", u.Id, subId).Scan()
	if err == sql.ErrNoRows {
		WriteJSONError(w, http.StatusNotFound, "Subscription doesn't exist")
		return
	}
	var article Article
	err = DB.QueryRow("select id from articles where subscription_id = ? and id = ?", subId, artId).Scan(&article.Id)
	if err == sql.ErrNoRows {
		WriteJSONError(w, http.StatusNotFound, "Article doesn't exist")
		return
	}
	var alreadyRead bool
	err = DB.QueryRow("select id from user_read_articles where article_id = ? and user_id = ?", artId, u.Id).Scan()
	if err == sql.ErrNoRows {
		alreadyRead = false
	} else {
		alreadyRead = true
	}
	switch status {
	case "read":
		if alreadyRead != true {
			DB.Exec("insert into user_read_articles values (null, ?, ?, ?)", u.Id, artId, false)
		}
	case "unread":
		if alreadyRead == true {
			DB.Exec("delete from user_read_articles where user_id = ? and article_id = ?", u.Id, artId)
		}
	default:
		WriteJSONError(w, http.StatusBadRequest, "Valid values for 'status' are 'read' and 'unread'")
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
status      = string || read | unread | all [all]
limit       = int    ||                     [100]
order       = string || asc | desc          [desc]
[since]     = string ||
[since_id]  = int    ||
[before]    = string ||
[before_id] = int    ||
*/
