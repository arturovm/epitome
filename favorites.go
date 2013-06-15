package main

import (
	"net/http"
)

type Favorite struct {
	Id string
	Url string
	Name string
	Author string
	Published int
	Body string
}

func PostFavorites(w http.ResponseWriter, req *http.Request) {
}

func GetFavorites(w http.ResponseWriter, req *http.Request) {
}

func DeleteFavorite(w http.ResponseWriter, req *http.Request) {
}
