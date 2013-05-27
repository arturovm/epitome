package main

import (
	"net/http"
)

type Article struct {
	Id     string
	Url    string
	Name   string
	Parent *Subscription
}

func GetAllArticles(w http.ResponseWriter, req *http.Request) {
}

func GetArticles(w http.ResponseWriter, req *http.Request) {
}

func PutArticle(w http.ResponseWriter, req *http.Request) {
}
