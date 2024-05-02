package services

import (
	"encoding/json"
	"io"
	"net/http"
)

type InMemoryService struct {
	db map[string]string
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		db: map[string]string{},
	}
}

func (ims *InMemoryService) Fetch(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, ok := ims.db[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	inMemoryPayload := InMemoryPayload{
		Key: key,
		Value: value,
	}

	responseJson, err := json.Marshal(inMemoryPayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

func (ims *InMemoryService) Store(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var inMemoryPayload InMemoryPayload
	if err  := json.Unmarshal(body, &inMemoryPayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	ims.db[inMemoryPayload.Key] = inMemoryPayload.Value

	responseJson, err := json.Marshal(inMemoryPayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
