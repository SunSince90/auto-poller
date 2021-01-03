// +build ignore

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	poller "github.com/SunSince90/website-poller"
	"gopkg.in/yaml.v2"
)

// This example shows how to start multiple pollers and manage their life
// cycles correctly.

func main() {
	// -- Load the file
	filePath := "./path/to/pages.yaml"

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("error while reading file:", err)
		os.Exit(1)
	}

	// -- Set up the pollers
	var pages []poller.Page
	var pollers []poller.Poller
	err = yaml.Unmarshal(yamlFile, &pages)
	if err != nil {
		fmt.Println("error while unmarshalling file:", err)
		os.Exit(1)
	}

	for _, page := range pages {
		p, err := poller.New(&page)
		if err != nil {
			fmt.Println("a page contains errors:", err, ", skipping...")
			continue
		}

		p.SetHandlerFunc(handleResponse)
		pollers = append(pollers, p)
	}

	// -- Start the pollers
	ctx, canc := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(len(pollers))

	for _, p := range pollers {

		go func(currentPoller poller.Poller) {
			defer wg.Done()
			currentPoller.Start(ctx, true)
		}(p)

	}

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
	canc()    // Cancel all the pollers
	wg.Wait() // Wait for all the pollers to finish before exiting!
	fmt.Println("goodbye!")
}

func handleResponse(id string, resp *http.Response, err error) {
	// Scrape the website...
}
