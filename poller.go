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
	defaultFrequency         int    = 30
	minFrequency             int    = 5
	minOffset                int    = 5
	minRandomFrequency       int    = 15
	defaultOffsetRange       int    = 10
	defaultHTTPClientTimeout int    = 20
	userAgentHeaderKey       string = "User-Agent"
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
	// GetID returns the ID of this poller. If the `Page` struct provided
	// to `New` contained a non-empty `ID`, then this returns the same ID as
	// the one contained in there, otherwise it returns a randomly generated
	// one.
	GetID() string
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
	if now {
		p.poll(ctx)
	}

	if !p.randTick {
		p.startFixed(ctx)
	} else {
		p.startRandom(ctx)
	}
}

func (p *pagePoller) startFixed(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.ticks) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go p.poll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (p *pagePoller) startRandom(ctx context.Context) {
	ticker := time.NewTimer(time.Duration(p.ticks) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go p.poll(ctx)
			next := nextRandomTick(p.ticks-p.offsetRange, p.ticks+p.offsetRange)
			ticker.Reset(time.Duration(next) * time.Second)
		case <-ctx.Done():
			return
		}
	}
}

func (p *pagePoller) poll(ctx context.Context) {
	// -- Get the user agent for this request,
	// and get the one for the next request
	userAgent, index := getNextUA(p.id, p.userAgents, p.randUa, p.lastUAIndex)
	p.lastUAIndex = index

	// -- Clone the request
	req := p.request.Clone(ctx)
	if len(userAgent) > 0 {
		req.Header.Set(userAgentHeaderKey, userAgent)
	}

	resp, err := p.httpClient.Do(req)

	// -- Pass response and error to the response handler func
	if p.HandlerFunc != nil {
		p.HandlerFunc(p.id, resp, err)
		return
	}
}

// SetHandlerFunc sets the function that will be called when a poll has
// finished
func (p *pagePoller) SetHandlerFunc(f HandlerFunc) {
	p.HandlerFunc = f
}

// GetID returns the ID of this poller. If the `Page` struct provided
// to `New` contained a non-empty `ID`, then this returns the same ID as
// the one contained in there, otherwise it returns a randomly generated
// one.
func (p *pagePoller) GetID() string {
	return p.id
}
