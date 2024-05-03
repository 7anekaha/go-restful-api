package services

import (
	"encoding/json"
	"fmt"
	"time"
)

type MongoRequestPayload struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

func (m *MongoRequestPayload) UnmarshalJSON(data []byte) error {
	type Alias MongoRequestPayload
	aux := &struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse StartDate
	startDate, err := time.Parse("2006-01-02", aux.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date. Format: 2024-12-31")
	}

	// Parse EndDate
	endDate, err := time.Parse("2006-01-02", aux.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date. Format: 2024-12-31")
	}

	// Assign parsed values to m fields
	m.StartDate = startDate
	m.EndDate = endDate
	m.MinCount = aux.MinCount
	m.MaxCount = aux.MaxCount

	return nil
}


type MongoResponsePayload struct {
	Code    uint8                `json:"code"`
	Msg     string               `json:"msg"`
	Records []MongoRecordPayload `json:"records"`
}

type MongoRecordPayload struct {
	Key        string `json:"key"`
	CreatedAt  time.Time `json:"createdAt"`
	TotalCount int    `json:"totalCount"`
}

type MongoRecord struct {
	Key       string    `json:"key" bson:"key"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Count     []int     `json:"count" bson:"count"`
}

type InMemoryPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
