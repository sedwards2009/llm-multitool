package main

type SessionOverview struct {
	SessionSummaries []SessionSummary `json:"sessionSummaries"`
}

type SessionSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Root struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Prompt    string     `json:"prompt"`
	Responses []Response `json:"responses"`
}

type Response struct {
	Text string `json:"text"`
}
