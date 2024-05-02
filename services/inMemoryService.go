package services

import "net/http"

type InMemoryService struct {
	db map[string]string
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		db: map[string]string{},
	}
}

func (ims *InMemoryService) Fetch(w http.ResponseWriter, r *http.Request) {
	// fetch data from memory
}

func (ims *InMemoryService) Store(w http.ResponseWriter, r *http.Request) {
	// store data in memory
}
