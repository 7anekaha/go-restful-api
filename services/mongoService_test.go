package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	// "log"
	// "os"
	"testing"

	"github.com/7anekaha/go-restful-api-week-17/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collection = "records"
	db         = "task17-db"
)

var mc *mongo.Client
var mongodbContainer *mongodb.MongoDBContainer

func TestMain(m *testing.M) {
	ctx := context.Background()
	populateDB(ctx)


	code := m.Run()

	if mc != nil {
		if err := mc.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
	mongodbContainer.Terminate(ctx)
	os.Exit(code)
}

func populateDB(ctx context.Context) {
	
	var err error
	mongodbContainer, err = mongodb.RunContainer(
		ctx,
		testcontainers.WithImage("mongo:6"),
		testcontainers.WithEnv(map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "username",
			"MONGO_INITDB_ROOT_PASSWORD": "supersecretpassword",
			"COLLECTION":                 collection,
			"DB":                         db,
		}),
	)
	if err != nil {
		log.Panic(err)
	}
	

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Panic(err)
	}

	// connect to mongo and populate the database
	mc, err = mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Panic(err)
	}

	// populate mongodb
	createdAtStr1 := "2017-01-27T08:19:14.135Z"
	createdAtStr2 := "2017-01-27T13:22:10.421Z"
	createdAtStr3 := "2017-01-28T01:22:14.398Z"
	createdAt1, _ := time.Parse(time.RFC3339, createdAtStr1)
	createdAt2, _ := time.Parse(time.RFC3339, createdAtStr2)
	createdAt3, _ := time.Parse(time.RFC3339, createdAtStr3)

	records := []interface{}{
		services.MongoRecord{
			Key:       "TAKwGc6Jr4i8Z487",
			CreatedAt: createdAt1,
			Count:     []int{500, 400, 450, 550, 300, 150, 350},
		},
		services.MongoRecord{
			Key:       "NAeQ8eX7e5TEg70H",
			CreatedAt: createdAt2,
			Count:     []int{540, 400, 450, 550, 300, 160, 350},
		},
		services.MongoRecord{
			Key:       "cCddT2RPqWmUI4Nf",
			CreatedAt: createdAt3,
			Count:     []int{120, 400, 450, 660, 500, 770, 250},
		},
	}

	res, err := mc.Database(db).Collection(collection).InsertMany(ctx, records)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.InsertedIDs)
	

}

func TestMongoService_Fetch(t *testing.T) {

	ms := services.NewMongoServiceForTest(mc)

	payload := map[string]any{
		"startDate": "2016-01-02",
		"endDate":   "2018-01-02",
		"minCount":  2800,
		"maxCount":  3200,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/mongodb", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ms.Fetch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var responsePayload services.MongoResponsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &responsePayload); err != nil {
		t.Error(err)
	}

	if responsePayload.Code != uint8(services.Success) {
		t.Errorf("Expected code 0, got %d", responsePayload.Code)
	}

	if responsePayload.Msg != "Success" {
		t.Errorf("Expected Success, got %s", responsePayload.Msg)
	}

	if len(responsePayload.Records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(responsePayload.Records))
	}
}


func TestMongoService_Fetch_0_records_found(t *testing.T) {

	ms := services.NewMongoServiceForTest(mc)

	payload := map[string]any{
		"startDate": "2016-01-02",
		"endDate":   "2018-01-02",
		"minCount":  4000,
		"maxCount":  5000,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/mongodb", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ms.Fetch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var responsePayload services.MongoResponsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &responsePayload); err != nil {
		t.Error(err)
	}

	if responsePayload.Code != uint8(services.Success) {
		t.Errorf("Expected code 0, got %d", responsePayload.Code)
	}

	if responsePayload.Msg != "Success" {
		t.Errorf("Expected Success, got %s", responsePayload.Msg)
	}

	if len(responsePayload.Records) != 0 {
		t.Errorf("Expected 1 record, got %d", len(responsePayload.Records))
	}
}

func TestMongoService_Fetch_bad_structure(t *testing.T) {

	ms := services.NewMongoServiceForTest(mc)

	payload := map[string]any{
		"endDate":   "2018-01-02",
		"minCount":  4000,
		"maxCount":  5000,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/mongodb", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ms.Fetch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", rr.Code)
	}

	var responsePayload services.MongoResponsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &responsePayload); err != nil {
		t.Error(err)
	}

	if responsePayload.Code != uint8(services.BadRequest) {
		t.Errorf("Expected code 2, got %d", responsePayload.Code)
	}

	if responsePayload.Msg != "Structure not valid" {
		t.Errorf("Expected Structure not valid, got %s", responsePayload.Msg)
	}

	if len(responsePayload.Records) != 0 {
		t.Errorf("Expected 0 record, got %d", len(responsePayload.Records))
	}
}

func TestMongoService_Fetch_minCount_is_greater_than_macCount(t *testing.T) {

	ms := services.NewMongoServiceForTest(mc)

	payload := map[string]any{
		"startDate": "2016-01-02",
		"endDate":   "2018-01-02",
		"minCount":  2000,
		"maxCount":  100,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/mongodb", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ms.Fetch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", rr.Code)
	}

	var responsePayload services.MongoResponsePayload
	if err := json.Unmarshal(rr.Body.Bytes(), &responsePayload); err != nil {
		t.Error(err)
	}

	if responsePayload.Code != uint8(services.BadRequest) {
		t.Errorf("Expected code 2, got %d", responsePayload.Code)
	}

	if responsePayload.Msg != "Start date should be before end date and minCount should be less than maxCount" {
		t.Errorf("Expected Start date should be before end date and minCount should be less than maxCount, got %s", responsePayload.Msg)
	}

	if len(responsePayload.Records) != 0 {
		t.Errorf("Expected 0 record, got %d", len(responsePayload.Records))
	}
}
