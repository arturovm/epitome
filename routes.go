package main

import (
	"net/http"
	"github.com/bmizerany/pat"
)

func RegisterRoutes() *pat.PatternServeMux {
	m := pat.New()
	// API
	// Auth
	m.Get("/api/auth", http.HandlerFunc(GetAuth))
	// Subscriptions
	m.Post("/api/subscriptions", http.HandlerFunc(PostSubscription))
	m.Get("/api/subscriptions", http.HandlerFunc(GetSubscriptions))
	/*
	* Regarding GET /subscriptions:
	* Should this list the user's subscriptions? Or return all his articles (with proper filters and such)?
	* The RESTful thing to do seems to be the first option. But then, what endpoint should return all of the user's articles?
	*/
	m.Del("/api/subscriptions/:id", http.HandlerFunc(DeleteSubscription))
	// Articles
	m.Get("/api/subscriptions/:id", http.HandlerFunc(GetArticles))
	m.Put("/api/subscriptions/:subid/:artid", http.HandlerFunc(PutArticle))
	// Favorites
	m.Post("/api/favorites", http.HandlerFunc(PostFavorites))
	m.Get("/api/favorites", http.HandlerFunc(GetFavorites))
	m.Del("/api/favorites/:id", http.HandlerFunc(DeleteFavorite))
	// First time setup
	m.Get("/api/setup", http.HandlerFunc(GetSetup)) // TODO: Is this really the best URL for this?
	return m
}
