package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/7anekaha/go-restful-api-week-17/services"
)


func TestInMemoryService_Store(t *testing.T) {
	ims := services.NewInMemoryService()

	payload := services.InMemoryPayload{
		Key: "key-test",
		Value: "value-test",
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/in-memory", bytes.NewBuffer(payloadJson))
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ims.Store)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", rr.Code)
	}

	if len(ims.DB) != 1 {
		t.Errorf("Expected 1 record, got %d", len(ims.DB))
	}

	if ims.DB["key-test"] != "value-test" {
		t.Errorf("Expected value-test, got %s", ims.DB["key-test"])
	}
}

func TestInMemoryService_Fetch(t *testing.T) {

	ims := services.NewInMemoryService()

	ims.DB["test-key"] = "test-value"

	req, err := http.NewRequest(http.MethodGet, "/in-memory?key=test-key", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ims.Fetch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var responsePayload services.InMemoryPayload
	if err := json.NewDecoder(rr.Body).Decode(&responsePayload); err != nil {
		t.Error(err)
	}
	
	if responsePayload.Key != "test-key" {
		t.Errorf("Expected test-key, got %s", responsePayload.Key)
	}

	if responsePayload.Value != "test-value" {
		t.Errorf("Expected test-value, got %s", responsePayload.Value)
	}
}
