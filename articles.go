package main

type Article struct {
	Id string
	Url string
	Name string
	Parent *Subscription
)
