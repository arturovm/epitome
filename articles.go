package main

import (
	"net/http"
	"time"
)

type ArticleContent struct {
	Type    string
	Content string
}

type Article struct {
	Id             string
	SubscriptionId int
	Url            string
	Title          string
	Author         string
	Published      time.Time
	Body           ArticleContent
	Summary        ArticleContent
	Read           bool
}

func GetAllArticles(w http.ResponseWriter, req *http.Request) {
}

func GetArticles(w http.ResponseWriter, req *http.Request) {
}

func PutArticle(w http.ResponseWriter, req *http.Request) {
}
