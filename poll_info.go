package main

import (
	"net/http"
	"sync"
)

type pollInfo struct {
	*WebsitePage
	remaining uint
	f         CallBack

	sync.Mutex
}

// CallBack is a pointer to a function that must be executed after
// an http call is performed
type CallBack func(*http.Response)
