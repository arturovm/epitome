package main

import (
	"net/http"
	"github.com/bmizerany/pat"
)

func RegisterRoutes() *pat.PatternServeMux {
	m := pat.New()
	// Auth
	m.Post("/apps", http.HandlerFunc(PostApp))
	m.Get("/apps/:id", http.HandlerFunc(GetApp))
	m.Put("/apps/:id", http.HandlerFunc(PutApp))
	// Subscriptions
	m.Post("/subscriptions", http.HandlerFunc(PostSubscription))
	m.Get("/subscriptions", http.HandlerFunc(GetSubscriptions))
	m.Delete("/subscriptions/:id", http.HandlerFunc(DeleteSubscription))
	// Articles
	m.Get("/subscriptions/:id", http.HandlerFunc(GetArticles))
	m.Put("/subscriptions/:subid/:artid", http.HandlerFunc(PutArticle))
	// Favorites
	m.Post("/favorites", http.HandlerFunc(PostFavorites))
	m.Delete("/favorites/:id", http.HandlerFunc(deleteFavorite))
	return m
}
