package main

import (
	"log"
	"net/http"
)

type Article struct {
	Id        string
	Url       string
	Name      string
	Author    string
	Published int
	Parent    *Subscription
	Body      string
	Read      bool
}

func GetAllArticles(w http.ResponseWriter, req *http.Request) {
}

func GetArticles(w http.ResponseWriter, req *http.Request) {
}

func PutArticle(w http.ResponseWriter, req *http.Request) {
}

func DownloadArticles() {
	log.Println("Downloading")
}
