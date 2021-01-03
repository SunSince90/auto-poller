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

// This is a simple example that defines some custom polling options.

func main() {
	id := "poll-user"
	offsetRange := 5
	page := &poller.Page{
		ID:  &id,
		URL: "https://api.github.com/users/sunsince90",
		Headers: map[string][]string{
			"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Accept-Encoding": {"gzip, deflate, br"},
			"Accept-Language": {"en-US,it-IT;q=0.8,it;q=0.5,en;q=0.3"},
			"Cache-Control":   {"no-cache"},
		},
		UserAgentOptions: &poller.UserAgentOptions{
			UserAgents: []string{
				"USER-AGENT-1",
				"USER-AGENT-2",
				"USER-AGENT-3",
				// ...
			},
			RandomUA: false, // rotate them at each request
		},
		PollOptions: &poller.PollOptions{
			Frequency:       20,           // poll every 20 seconds
			RandomFrequency: true,         // mimick user behavior: don't make requests at a fixed time
			OffsetRange:     &offsetRange, // offsetRange is 5, so the range of requests is [15 - 25] seconds
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
