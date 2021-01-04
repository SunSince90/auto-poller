package websitepoller

import "context"

// Poller is in charge of polling a website and providing results to a
// function that will handle the result
type Poller interface {
	// Start polling
	Start(ctx context.Context, now bool)
	// SetHandlerFunc sets the function that will be called when a poll has
	// finished
	SetHandlerFunc(HandlerFunc)
	// GetID returns the ID of this poller. If the `Page` struct provided
	// to `New` contained a non-empty `ID`, then this returns the same ID as
	// the one contained in there, otherwise it returns a randomly generated
	// one.
	GetID() string
}
