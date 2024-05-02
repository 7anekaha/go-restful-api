package services

import "net/http"

type Fetcher interface {
	Fetch(w http.ResponseWriter, r *http.Request)
}

type Storer interface {
	Store(w http.ResponseWriter, r *http.Request)
}
