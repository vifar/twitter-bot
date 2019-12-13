package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	log "github.com/sirupsen/logrus"
)

const dots = '⋮'

var count = 1

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.Info("Retrieving Keys.......")

	client := auth()

	// Get current Date & Time
	localTime, _ := time.LoadLocation("UTC")
	now := time.Now().In(localTime)
	log.Infoln("Today is: ", localTime, " & the time is: ", now)

	// Years to do calculations with
	currentYear := now.Year()
	nextYear := currentYear + 1
	decadeEnd := int(math.Round(float64(currentYear)/10) * 10)

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				calcYearCompleted(currentYear, nextYear, client)
				calcDecadeCompleted(currentYear, decadeEnd, client)
			case <-quit:
				log.Info("Retrieving Keys.......")
				return
			}
		}
	}()
	defer ticker.Stop()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

func sendTweet(status string, percent int) {
	tweet := fmt.Sprintf("%s - %d%", status, percent)
	//Send a Tweet
	_, _, err := client.Statuses.Update(tweet, nil)
	if err != nil {
		log.Error(err)
	}
}

func calcYearCompleted(currentYear int, nextYear int, client *twitter.Client) {

	// sendTweet(client)

}

func calcDecadeCompleted(currentYear int, decadeEnd int, client *twitter.Client) {

	// sendTweet(client)

}
