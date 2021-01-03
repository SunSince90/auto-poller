// +build ignore

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	poller "github.com/SunSince90/website-poller"
)

// This is a simple example that just logs the result of each request.
// No parameteres are defined, so, it wil poll every 30 seconds with no
// user agents.

func main() {
	id := "poll-user"
	page := &poller.Page{
		ID:  &id,
		URL: "https://api.github.com/users/sunsince90",
		Headers: map[string][]string{
			"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Accept-Encoding": {"gzip, deflate, br"},
			"Accept-Language": {"en-US,it-IT;q=0.8,it;q=0.5,en;q=0.3"},
			"Cache-Control":   {"no-cache"},
		},
	}

	p, err := poller.New(page)
	if err != nil {
		fmt.Println("error occurred:", err)
		os.Exit(1)
	}
	p.SetHandlerFunc(handleResponse)

	ctx, canc := context.WithCancel(context.Background())
	exitChan := make(chan struct{})

	go func() {
		p.Start(ctx, true)
		close(exitChan)
	}()

	// -- Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	fmt.Println("exit requested")
	canc()
	<-exitChan // Wait for the poller goroutine to return
	fmt.Println("goodbye!")
}

func handleResponse(id string, resp *http.Response, err error) {
	fmt.Println("request with id", id, "returned status", resp.Status)
}
