package services

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)
type MongoService struct {
	db *mongo.Client
}

func NewMongoService() *MongoService {
	return &MongoService{
		db : nil,
	}
}

func (ms *MongoService) Fetch(w http.ResponseWriter, r *http.Request) {
	// fetch data from mongodb
}
