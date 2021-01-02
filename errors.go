package websitepoller

import "errors"

var (
	// ErrUnrecognizedHTTPMethod means that the http method is not recognized
	ErrUnrecognizedHTTPMethod = errors.New("unrecognized http method")
	// ErrURLNoScheme means that the url has no http:// or https://
	ErrURLNoScheme = errors.New("url does not have a scheme")
	// ErrInvalidFrequency means that the polling frequency is invalid, i.e.
	// if it could not be parsed by time.ParseDuration
	ErrInvalidFrequency = errors.New("invalid polling frequency")
	// ErrUnsupportedFrequency means that the polling frequency is lower
	// than five second
	ErrUnsupportedFrequency = errors.New("frequency cannot be lower than five second")
	// ErrInvalidRandRange is thrown when the range could not be parsed by
	// time.ParseDuration or is greater or equal than frequency
	ErrInvalidRandRange = errors.New("invalid range")
)
