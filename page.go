package websitepoller

import "net/http"

// Page to poll
type Page struct {
	// ID is a short name that will be used by the logs to recognize when
	// operations are performed on this page.
	ID *string `yaml:"id,omitempty"`
	// URL to poll
	URL string `yaml:"url"`
	// Method of the request, e.g.: GET
	Method *string `yaml:"method,omitempty"`
	// Headers to send with the request
	Headers map[string]string `yaml:"headers,omitempty"`
	// UserAgentOptions contains options about the
	// user agent
	*UserAgentOptions `yaml:"userAgentOptions,omitempty"`
	// PollOptions contains options about polling
	*PollOptions `yaml:"pollOptions,omitempty"`

	// TODO: support cookies?
	// TODO: support body
}

// UserAgentOptions contains options about the user agent
type UserAgentOptions struct {
	// UserAgents is a list of user agents that should be used for this page
	UserAgents []string `yaml:"userAgents,omitempty"`
	// RandomUA specifies whether the user agent for the next request should
	// be chosen randomly or should be rotated. This has no effect if
	// UserAgents is empty or only has one element. Leave this false if you
	// have few user agents.
	RandomUA bool `yaml:"randomUA"`
}

// PollOptions contains options about polling
type PollOptions struct {
	// Frequency of polling in seconds
	Frequency int `yaml:"frequency"`
	// Randomize specifies whether the next poll should be at a random time or
	// at a fixed time
	RandomFrequency bool `yaml:"randomFrequency"`
	// OffsetRange specifies the range for choosing the next random time.
	// For example, if Frequency is 30 and OffsetRange is 10, then each next
	// poll will be performed at a random time in the [20, 40] seconds range,
	// i.e. 27 seconds.
	OffsetRange *int `yaml:"offsetRange,omitempty"`
}

// HandlerFunc represents a function that will handle the response returned
// by the polling.
type HandlerFunc func(string, *http.Response, error)
