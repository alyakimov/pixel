package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Redirect",
		"GET",
		"/r",
		Redirect,
	},
	Route{
		"Index",
		"GET",
		"/{(.*)}",
		Index,
	},
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
}
