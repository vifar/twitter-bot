package main

import (
	"fmt"
	"math"
	"net/http"
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

var yearProgress = 0
var decadeProgress = 0

func main() {

	log.SetFormatter(&log.TextFormatter{})
	log.Info("Retrieving Keys.......")

	client := auth()

	// Get current Date & Time
	localTime, _ := time.LoadLocation("UTC")
	now := time.Now().In(localTime)
	log.Infoln("The time is: ", now)

	// Years to do calculations with
	nextYear := now.Year() + 1
	decadeEnd := int(math.Round(float64(now.Year())/10) * 10)
	// decadeBeginning := decadeEnd - 10

	calcYearCompleted(time.Now().In(localTime), nextYear, client)

	ticker := time.NewTicker(time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Info("Checking time passed.......")
				calcYearCompleted(time.Now().In(localTime), nextYear, client)
				calcDecadeCompleted(time.Now().In(localTime), decadeEnd, client)
			case <-quit:
				log.Info("Quiting.......")
				return
			}
		}
	}()
	defer ticker.Stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, nil)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println(<-ch)

}

func sendTweet(status string, percent int, now time.Time, client *twitter.Client) {

	status = fmt.Sprintf("Year %[1]d\n\n%[2]s - %[3]d%%", now.Year(), status, percent)

	log.Info("Sending teweet.......")

	// Send a Tweet
	_, _, err := client.Statuses.Update(status, nil)
	if err != nil {
		log.Error(err)
	}

}

func calcYearCompleted(now time.Time, nextYear int, client *twitter.Client) {

	log.Info("Calculating year completed.......")

	difference := now.Sub(time.Date(nextYear, time.January, 1, 0, 0, 0, 0, time.UTC))

	var percent int
	if (now.Year() % 4) == 0 {
		percent = int(((daysInLeapYear - (difference.Hours() / -24)) / daysInLeapYear) * 100)
	}
	percent = int(((daysInYear - (difference.Hours() / -24)) / daysInYear) * 100)

	log.Info("Percent: ", percent, " & Year Progress: ", yearProgress)

	if percent > yearProgress {

		log.Info("Composting teweet.......")

		yearProgress = percent
		var status string
		for i := 0; i <= (100 * .5); i = i + 4 {
			if float64(i) <= (float64(percent) * .5) {
				status += completed
			} else {
				status += notCompleted
			}
		}
		sendTweet(status, percent, now, client)
	}

}

func calcDecadeCompleted(now time.Time, decadeEnd int, client *twitter.Client) {

	//difference := now.Sub(time.Date(decadeEnd, time.January, 1, 0, 0, 0, 0, time.UTC))

	// sendTweet(client)

}
