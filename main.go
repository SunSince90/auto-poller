package autopoller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, canc := context.WithCancel(context.Background())
	defer canc()

	exitChan := make(chan struct{})

	p := New(10, false)
	var freq int = 10
	p.AddPage(&WebsitePage{
		ID:  "euronics-digital",
		URL: "https://www.euronics.it/console/sony-computer/playstation-5-digital-edition/eProd202008907/",
		PollSettings: PollSettings{
			Type:      FixedPolling,
			Frequency: &freq,
		},
		UserAgents: []string{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
	}, do)
	p.AddPage(&WebsitePage{
		ID:  "euronics-standard",
		URL: "https://www.euronics.it/console/sony-computer/playstation-5/eProd202008906/",
		PollSettings: PollSettings{
			Type:      FixedPolling,
			Frequency: &freq,
		},
		UserAgents: []string{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
	}, do)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	go p.Start(ctx, exitChan)

	<-signalChan
	fmt.Println("os.Interrupt - shutting down...")
	canc()

	// PERFORM GRACEFUL SHUTDOWN HERE
	<-exitChan // Wait for euronics to finish

	fmt.Println("exiting for real")
	os.Exit(0)
}

func do(w *WebsitePage, r *http.Response, e error) {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-13-02 15:04:05"
	customFormatter.FullTimestamp = true
	l := log.WithFields(log.Fields{"id": w.ID, "resp": r.Status})
	l.Logger.SetFormatter(customFormatter)
	l.Info("got response")
	defer r.Body.Close()
}
