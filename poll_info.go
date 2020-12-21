package autopoller

import (
	"net/http"
	"sync"
)

type pollInfo struct {
	*WebsitePage
	remaining int
	f         CallBack

	lock sync.Mutex
}

// CallBack is a pointer to a function that must be executed after
// an http call is performed
type CallBack func(*WebsitePage, *http.Response, error)

func newPollInfo(page *WebsitePage, f CallBack) *pollInfo {
	return &pollInfo{
		WebsitePage: page,
		f:           f,
		remaining:   nextRandom(10, 60),
	}
}

func (i *pollInfo) shouldGo() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.remaining == 0 {
		if i.WebsitePage.PollSettings.Type == FixedPolling {
			i.remaining = *i.WebsitePage.Frequency
		} else {
			i.remaining = nextRandom(*i.WebsitePage.RandMin, *i.WebsitePage.RandMax)
		}

		return true
	}

	i.remaining--
	return false
}
