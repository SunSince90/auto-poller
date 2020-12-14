package main

// PollType represents the polling type (fixed or randomized)
type PollType string

const (
	// FixedPolling is a constant polling
	FixedPolling string = "fixed"
	// RandPolling means that the next poll will always be a random
	// value
	RandPolling string = "random"
)

// WebsitePage contains information of the page to poll
type WebsitePage struct {
	// ID of this page
	ID string `json:"id"`
	// URL to poll
	URL string `json:"url"`
	// UserAgents to use
	UserAgents   []string `json:"userAgents"`
	PollSettings `json:"pollSettings"`
}

// PollSettings contains settings about polling
type PollSettings struct {
	// Type of polling
	Type PollType `json:"type"`
	// Frequency of polling, in seconds
	Frequency *uint `json:"frequency"`
	// RandMin is the minimum value that can be extracted
	// when random polling
	RandMin *uint `json:"randMin"`
	// RandMin is the maximum value that can be extracted
	// when random polling
	RandMax *uint `json:"randMax"`
}
