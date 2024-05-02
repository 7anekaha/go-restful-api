package services

type MongoRequestPayload struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	MinCount int `json:"minCount"`
	MaxCount int `json:"maxCount"`
}

type MongoResponsePayload struct {
	Code uint8 `json:"code"`
	Msg string `json:"msg"`
	Records []Record `json:"records"`
}

type Record struct {
	Key string `json:"key"`
	CreatedAt string `json:"createdAt"`
	TotalCount int `json:"totalCount"`
}

type InMemoryPayload struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
