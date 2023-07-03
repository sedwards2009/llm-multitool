package data

import (
	"encoding/json"
	"testing"
)

func TestResponseStatusMarshall(t *testing.T) {
	response := &Response{
		Status:            ResponseStatus_Pending,
		CreationTimestamp: "",
		Prompt:            "A prompt",
		Text:              "",
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

	if response2.Status != ResponseStatus_Pending {
		t.Errorf("Round trip ResponseStatus is wrong.")
	}
}
