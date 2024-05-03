package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName string = os.Getenv("DB")
var collectionName string = os.Getenv("COLLECTION")
var uri string = os.Getenv("MONGO_URI")

type ResponseCode uint8

const (
	Success ResponseCode = iota
	InternalServerError
	BadRequest
)

type MongoService struct {
	db *mongo.Client
}

func NewMongoService(ctx context.Context) *MongoService {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	// Check the connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to MongoDB")

	return &MongoService{
		db: mongoClient,
	}
}

func (ms *MongoService) Fetch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	responsePayload := MongoResponsePayload{
		Code:    uint8(Success),
		Msg:     "Success",
		Records: nil,
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var requestPayload MongoRequestPayload
	if err := json.Unmarshal(body, &requestPayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responsePayload.Code = uint8(BadRequest)
		responsePayload.Msg = "Structure not valid"

		responseJson, err := json.Marshal(responsePayload)
		if err != nil {
			return
		}
		w.Write(responseJson)
		return
	}

	if ok := validatePayload(requestPayload); !ok {
		w.WriteHeader(http.StatusBadRequest)
		responsePayload.Code = uint8(BadRequest)
		responsePayload.Msg = "Start date should be before end date and minCount should be less than maxCount"

		responseJson, err := json.Marshal(responsePayload)
		if err != nil {
			return
		}
		w.Write(responseJson)
		return
	}
	filter := bson.M{
		"createdAt": bson.M{"$gte": requestPayload.StartDate,
			"$lte": requestPayload.EndDate},
	}
	cursor, err := ms.db.Database(dbName).Collection(collectionName).Find(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var records []MongoRecordPayload
	// iterate over cursor and count the sum of the counts to check if it is between minCount and maxCount
	for cursor.Next(r.Context()) {
		var record MongoRecord
		if err := cursor.Decode(&record); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var sum int
		for _, count := range record.Count {
			sum += count
		}
		if sum >= requestPayload.MinCount && sum <= requestPayload.MaxCount {
			recordPayload := MongoRecordPayload{
				Key:        record.Key,
				CreatedAt:  record.CreatedAt,
				TotalCount: sum,
			}
			records = append(records, recordPayload)
		}
	}

	responsePayload.Records = records
	responseJson, err := json.Marshal(responsePayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(responseJson)
}

func validatePayload(payload MongoRequestPayload) bool {
	if payload.StartDate.After(payload.EndDate) || payload.MinCount >= payload.MaxCount {
		return false
	}
	return true
}
