package autopoller

// PollType represents the polling type (fixed or randomized)
type PollType string

const (
	// FixedPolling is a constant polling
	FixedPolling PollType = "fixed"
	// RandPolling means that the next poll will always be a random
	// value
	RandPolling PollType = "random"
)

// WebsitePage contains information of the page to poll
type WebsitePage struct {
	// ID of this page
	ID string `json:"id" yaml:"id"`
	// URL to poll
	URL string `json:"url" yaml:"url"`
	// NotifyURL is the url to notify externally
	NotifyURL *string `json:"notifyUrl" yaml:"notifyUrl,omitempty"`
	// UserAgents to use
	UserAgents []string `json:"userAgents" yaml:"userAgents,omitempty"`
	// PollSettings contains settings about polling
	PollSettings `json:"pollSettings" yaml:"pollSettings"`
}

// PollSettings contains settings about polling
type PollSettings struct {
	// Type of polling
	Type PollType `json:"type" yaml:"type"`
	// Frequency of polling, in seconds
	Frequency *int `json:"frequency" yaml:"frequency,omitempty"`
	// RandMin is the minimum value that can be extracted
	// when random polling
	RandMin *int `json:"randMin" yaml:"randMin,omitempty"`
	// RandMin is the maximum value that can be extracted
	// when random polling
	RandMax *int `json:"randMax" yaml:"randMax,omitempty"`
}
