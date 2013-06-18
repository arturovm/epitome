package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Subscription struct {
	Id   string
	Url  string
	Name string
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
	doc, err := gokogiri.ParseXml(body)
	defer doc.Free()
	if doc.Root().Name() != "rss" || err != nil {
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
					u.Host = removeTrail(rawurl)
					return u.String(), nil
				}
			}
		}
	} else {
		return res.Request.URL.String(), nil
	}
	return "", errors.New("URL not found")
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
	// Auth should go here, but well...
	if rawurl := req.FormValue("url"); rawurl != "" {
		rssUrl, err := findRSSURL(rawurl)
		if err != nil {
			// TODO: return error here
			log.Fatal(err)
		} else {
			log.Print(rssUrl)
			title, err := findRSSTitle(rssUrl)
			if err != nil {
				// TODO: return error here
				log.Fatal(err)
			}
			DB, _ := sql.Open("sqlite3", "db.db")
			DB.Exec("insert into subscriptions values (null, ?, ?)", rssUrl, title)
			DB.Close()
		}
	}
}

func GetSubscriptions(w http.ResponseWriter, req *http.Request) {
}

func DeleteSubscription(w http.ResponseWriter, req *http.Request) {
}
