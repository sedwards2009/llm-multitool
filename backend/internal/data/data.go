package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SessionOverview struct {
	SessionSummaries []*SessionSummary `json:"sessionSummaries"`
}

type SessionSummary struct {
	ID                string `json:"id"`
	CreationTimestamp string `json:"creationTimestamp"`
	Title             string `json:"title"`
}

type Root struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	ID                string         `json:"id"`
	CreationTimestamp string         `json:"creationTimestamp"`
	Title             string         `json:"title"`
	Prompt            string         `json:"prompt"`
	Responses         []*Response    `json:"responses"`
	ModelSettings     *ModelSettings `json:"modelSettings"`
}

type ModelSettings struct {
	ModelID string `json:"modelId"`
}

type ModelOverview struct {
	Models []*Model `json:"models"`
}

type Model struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	InternalModelID string
}

// -------------------------------------------------------------------------
type ResponseStatus uint8

const (
	ResponseStatus_Done ResponseStatus = iota + 1
	ResponseStatus_Pending
	ResponseStatus_Running
	ResponseStatus_Error
)

func (status ResponseStatus) String() string {
	switch status {
	case ResponseStatus_Pending:
		return "pending"
	case ResponseStatus_Done:
		return "done"
	case ResponseStatus_Running:
		return "running"
	case ResponseStatus_Error:
		return "error"
	default:
	}
	panic("Unknown ResponseStatus enum value.")
}

func (r ResponseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *ResponseStatus) UnmarshalJSON(data []byte) (err error) {
	var responseStatusString string
	if err := json.Unmarshal(data, &responseStatusString); err != nil {
		return err
	}

	if *r, err = ParseResponseStatus(responseStatusString); err != nil {
		return err
	}
	return nil
}

func ParseResponseStatus(statusString string) (ResponseStatus, error) {
	statusString = strings.TrimSpace(strings.ToLower(statusString))
	switch statusString {
	case "done":
		return ResponseStatus_Done, nil
	case "pending":
		return ResponseStatus_Pending, nil
	case "running":
		return ResponseStatus_Running, nil
	case "error":
		return ResponseStatus_Error, nil
	default:
		return ResponseStatus_Done, fmt.Errorf("%q is not a valid ResponseStatus value", statusString)
	}
}

//-------------------------------------------------------------------------

type Response struct {
	ID                string         `json:"id"`
	CreationTimestamp string         `json:"creationTimestamp"`
	Status            ResponseStatus `json:"status"`
	Prompt            string         `json:"prompt"`
	Text              string         `json:"text"`
}
