package main

import (
	"net/http"
)

type Subscription struct {
	Id string
	Url string
	Name string
}

func PostSubscription(w http.ResponseWriter, req *http.Request) {
}

func GetSubscriptions(w http.ResponseWriter, req *http.Request) {
}

func DeleteSubscription(w http.ResponseWriter, req *http.Request) {
