package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	log.Info("Retrieving Keys.......")

	client := auth()

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}

	getRandomWord()
	// Send a Tweet
	// _, _, err := client.Statuses.Update("just setting up my twttr", nil)
	// if err != nil {
	// 	log.Error(err)
	// }

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"@Future"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()

}
