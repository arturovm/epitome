package main

import (
	"log"
	"net/http"
)

func GetSetup(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(401)
	w.Write([]byte("You are not authorized to view this page"))
}

func GetStatic(w http.ResponseWriter, req *http.Request) {
	//w.Header()["content-type"] = http.DetectContentType(fileio.BytesFromFile("static/" + req.URL.Query().Get(":file") + pat.Tail("/:file", )))
	log.Fatal(req.URL.Path)
}
