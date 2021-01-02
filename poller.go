package websitepoller

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var (
	log zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// Poller is in charge of polling a website and providing results to a
// function that will handle the result
type Poller interface {
	// Start polling
	Start(ctx context.Context, now bool)
	// SetHandlerFunc sets the function that will be called when a poll has
	// finished
	SetHandlerFunc(HandlerFunc)
}

type pagePoller struct {
	// TODO: implement me
}

// New returns a new instance of the poller
func New(p *Page) (Poller, error) {
	// TODO: implement me
	return nil, nil
}

// Start polling
func (p *pagePoller) Start(ctx context.Context, now bool) {
	// TODO: implement me
}

// SetHandlerFunc sets the function that will be called when a poll has
// finished
func (p *pagePoller) SetHandlerFunc(f HandlerFunc) {
	// TODO: implement me
}
