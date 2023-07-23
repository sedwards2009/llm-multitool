package data

import (
	"encoding/json"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/data/role"
	"testing"
)

func TestResponseStatusMarshall(t *testing.T) {
	response := &Response{
		Status:            responsestatus.Pending,
		CreationTimestamp: "",
		Messages: []Message{
			{
				Role: role.User,
				Text: "A prompt",
			},
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Couldn't marshal Session object: %v", err)
		return
	}

	var response2 Response
	err = json.Unmarshal([]byte(jsonData), &response2)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
		return
	}

	if response2.Status != responsestatus.Pending {
		t.Errorf("Round trip ResponseStatus is wrong.")
	}
}
