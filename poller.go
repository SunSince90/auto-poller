package autopoller

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Poller performs continous polling to websites
type Poller interface {
	Start(ctx context.Context, exitChan chan struct{})
	AddUserAgents(uas []string)
	AddPage(page *WebsitePage, f CallBack)
	RemovePage(id string)
}

type poll struct {
	pages      map[string]*pollInfo
	globalUAs  []string
	httpClient *http.Client
	lock       sync.Mutex
}

// New returns a new instance of the poller
func New(timeout int, redir bool) Poller {
	to := time.Duration(timeout) * time.Second
	httpClient := &http.Client{
		Timeout: to,
	}

	if !redir {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &poll{
		pages:      map[string]*pollInfo{},
		globalUAs:  []string{},
		httpClient: httpClient,
	}
}

// AddUserAgents adds user agents for all pages, except those that
// set their owns
func (p *poll) AddUserAgents(uas []string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.globalUAs = append(p.globalUAs, uas...)
}

// Start polling
func (p *poll) Start(ctx context.Context, exitChan chan struct{}) {
	l := log.WithField("func", "poll.Start")
	l.Info("starting...")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			p.check(ctx)
		case <-ctx.Done():
			l.Info("stop requested")
			ticker.Stop()
			close(exitChan)
			return
		}
	}
}

// AddPage adds a new website to poll
func (p *poll) AddPage(page *WebsitePage, f CallBack) {
	l := log.WithFields(log.Fields{"func": "poll.AddPage", "id": page.ID})
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pages[page.ID] = newPollInfo(page, f)
	l.WithField("polling-in", p.pages[page.ID].remaining).Info("added")
}

// RemovePage removes website page
func (p *poll) RemovePage(id string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.pages, id)
}

func (p *poll) check(ctx context.Context) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, page := range p.pages {
		if page.shouldGo() {
			go p.do(ctx, page)
		}
	}
}

func (p *poll) do(ctx context.Context, pi *pollInfo) {
	l := log.WithFields(log.Fields{"func": "poll.do", "id": pi.ID})
	uas := p.globalUAs
	ua := ""

	// Get the user agent
	if len(pi.UserAgents) > 0 {
		if len(pi.UserAgents) == 0 {
			ua = pi.UserAgents[0]
		}
		uas = pi.UserAgents
	} else {
		if len(p.globalUAs) > 0 {
			n := nextRandom(0, len(uas)-1)
			ua = uas[n]
		}
	}

	l.Debug("using user agent", ua)

	resp, err := p.doRequest(ctx, pi.URL, ua)
	if err != nil {
		l.WithError(err).Error("error while doing request")
	}

	// Pass this to the callback
	go pi.f(pi.WebsitePage, resp, err)
}

func (p *poll) doRequest(ctx context.Context, url, ua string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	if len(ua) > 0 {
		req.Header.Add("User-Agent", ua)
	}

	resp, err = p.httpClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
			err = context.DeadlineExceeded
		}
	}

	return
}
