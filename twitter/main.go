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

const notCompleted = "░"
const completed = "█"
const daysInYear = 365
const daysInLeapYear = 366

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.Info("Retrieving Keys.......")

	client := auth()

	// Get current Date & Time
	localTime, _ := time.LoadLocation("UTC")
	now := time.Now().In(localTime)
	log.Infoln("Today is: ", localTime, " & the time is: ", now)

	// Years to do calculations with
	nextYear := now.Year() + 1
	decadeEnd := int(math.Round(float64(now.Year())/10) * 10)
	// decadeBeginning := decadeEnd - 10

	ticker := time.NewTicker(5 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				calcYearCompleted(time.Now().In(localTime), nextYear, client)
				calcDecadeCompleted(time.Now().In(localTime), decadeEnd, client)
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

func sendTweet(status string, percent int, now time.Time, client *twitter.Client) {

	status = fmt.Sprintf("Year of %[1]d\n\n%[2]s - %[3]d%%", now.Year(), status, percent)

	// Send a Tweet
	_, _, err := client.Statuses.Update(status, nil)
	if err != nil {
		log.Error(err)
	}

}

func calcYearCompleted(now time.Time, nextYear int, client *twitter.Client) {

	difference := now.Sub(time.Date(nextYear, time.January, 1, 0, 0, 0, 0, time.UTC))

	var percent int
	if (now.Year() % 4) == 0 {
		percent = int(((daysInLeapYear - (difference.Hours() / -24)) / daysInLeapYear) * 100)
	}
	percent = 25 //int(((daysInYear - (difference.Hours() / -24)) / daysInYear) * 100)

	var status string
	for i := 0; i <= 100; i = i + 4 {
		if i <= percent {
			status += completed
		} else {
			status += notCompleted
		}
	}
	sendTweet(status, percent, now, client)

}

func calcDecadeCompleted(now time.Time, decadeEnd int, client *twitter.Client) {

	//difference := now.Sub(time.Date(decadeEnd, time.January, 1, 0, 0, 0, 0, time.UTC))

	// sendTweet(client)

}
