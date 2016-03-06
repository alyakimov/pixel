package main

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(response, request)

		log.Printf(
			"%s\t%s\t%s\t(%s)\t%s",
			name,
			request.Method,
			request.RequestURI,
			request.RemoteAddr,
			time.Since(start),
		)
	})
}
