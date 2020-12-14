package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Poller performs continous polling to websites
type Poller interface {
	Start(ctx context.Context, exitChan chan struct{})
	AddPage(page *WebsitePage, f CallBack)
	RemovePage(id string)
}

type poll struct {
	pages map[string]*pollInfo

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
		httpClient: httpClient,
	}
}

// Start polling
func (p *poll) Start(ctx context.Context, exitChan chan struct{}) {
	l := log.WithField("func", "poll.Start")
	l.Info("starting...")

	ticker := time.NewTicker(time.Second)
	for {
		// Which one happens first?
		select {
		case <-ticker.C:
			p.check()
		case <-ctx.Done():
			l.Info("stop requested")
			ticker.Stop()
			exitChan <- struct{}{}
			return
		}
	}
}

// AddPage adds a new website to poll
func (p *poll) AddPage(page *WebsitePage, f CallBack) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pages[page.ID] = newPollInfo(page, f)
}

// RemovePage removes website page
func (p *poll) RemovePage(id string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.pages, id)
}

func (p *poll) check() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, page := range p.pages {
		if page.shouldGo() {
			// go poll
		}
	}
}
