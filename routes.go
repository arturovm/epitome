package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

func RegisterRoutes() *pat.PatternServeMux {
	m := pat.New()
	// API
	// Auth
	m.Get("/api/auth", http.HandlerFunc(GetAuth))
	// Users
	m.Post("/api/users", http.HandlerFunc(PostUser))
	m.Get("/api/users/:id", http.HandlerFunc(GetUser))
	m.Put("/api/users/:id", http.HandlerFunc(PutUsers))
	m.Del("/api/users/:id", http.HandlerFunc(DeleteUser))
	// Subscriptions
	m.Post("/api/subscriptions", http.HandlerFunc(PostSubscription))
	m.Get("/api/subscriptions", http.HandlerFunc(GetSubscriptions))
	m.Del("/api/subscriptions/:id", http.HandlerFunc(DeleteSubscription))
	// Articles
	m.Get("/api/subscriptions/articles", http.HandlerFunc(GetAllArticles))
	m.Get("/api/subscriptions/:id/articles", http.HandlerFunc(GetArticles))
	m.Put("/api/subscriptions/:subid/articles/:artid", http.HandlerFunc(PutArticle))
	// Favorites
	m.Post("/api/favorites", http.HandlerFunc(PostFavorites))
	m.Get("/api/favorites", http.HandlerFunc(GetFavorites))
	m.Del("/api/favorites/:id", http.HandlerFunc(DeleteFavorite))
	// First time setup
	m.Get("/api/setup", http.HandlerFunc(GetSetup)) // TODO: Is this really the best URL for this?
	return m
}
