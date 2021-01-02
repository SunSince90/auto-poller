package websitepoller

import (
	"context"
	"net/http"
	"os"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/rs/zerolog"
)

var (
	log zerolog.Logger
)

const (
	defaultFrequency         int = 30
	minFrequency             int = 5
	minOffset                int = 5
	minRandomFrequency       int = 15
	defaultOffsetRange       int = 10
	defaultHTTPClientTimeout int = 20
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
	id          string
	httpClient  *http.Client
	request     *http.Request
	userAgents  []string
	ticks       int
	randTick    bool
	offsetRange int
	lastUAIndex int
	randUa      bool
	HandlerFunc
}

// New returns a new instance of the poller
func New(p *Page) (Poller, error) {
	id := ""
	l := log.With().Str("func", "poller.New").Logger()

	if p.ID != nil && len(*p.ID) > 0 {
		id = *p.ID
	} else {
		id = randomdata.SillyName()
		l.Info().Msg("generating random name...")
	}
	l = l.With().Str("id", id).Logger()

	// -- Validation
	method, err := parseHTTPMethod(p.Method)
	if err != nil {
		return nil, err
	}
	parsedURL, err := parseURL(p.URL)
	if err != nil {
		return nil, err
	}

	// -- Set ups
	randomFrequency, ticks, offset := parsePollOptions(id, p.PollOptions)
	if randomFrequency {
		ticks = nextRandomTick(ticks-offset, ticks+offset)
	}

	randUA, userAgents := parseUserAgentOptions(id, p.UserAgentOptions)

	headers := http.Header{}
	if p.Headers == nil {
		l.Warn().Msg("no headers provided")
	} else {
		headers = p.Headers
		switch hlen := len(p.Headers); {
		case hlen == 0:
			l.Warn().Msg("no headers provided")
		case hlen < 3:
			l.Warn().Msg("few headers provided")
		}
	}

	httpClient := &http.Client{
		Timeout: time.Duration(defaultHTTPClientTimeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, err := http.NewRequestWithContext(context.Background(), method, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers

	// -- Complete and return
	return &pagePoller{
		id:          id,
		httpClient:  httpClient,
		request:     request,
		userAgents:  userAgents,
		ticks:       ticks,
		randTick:    randomFrequency,
		offsetRange: offset,
		lastUAIndex: -1,
		randUa:      randUA,
	}, nil
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
