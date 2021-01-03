// +build ignore

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	poller "github.com/SunSince90/website-poller"
	"gopkg.in/yaml.v2"
)

// This is a simple example loads pages to poll from a file.

func main() {
	filePath := "./path/to/pages.yaml"

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("error while reading file:", err)
		os.Exit(1)
	}

	var page poller.Page
	err = yaml.Unmarshal(yamlFile, &page)
	if err != nil {
		fmt.Println("error while unmarshalling file:", err)
		os.Exit(1)
	}
	p, err := poller.New(&page)
	if err != nil {
		fmt.Println("error while reading file:", err)
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
