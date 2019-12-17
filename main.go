package main

import (
	"encoding/json"
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
	calcDecadeCompleted(time.Now().In(localTime), decadeEnd, client)

	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Info("Checking time passed.......")
				// calcYearCompleted(time.Now().In(localTime), nextYear, client)
				calcDecadeCompleted(time.Now().In(localTime), decadeEnd, client)
			case <-quit:
				log.Info("Quiting.......")
				return
			}
		}
	}()
	defer ticker.Stop()

	// start new relic, and setting address to ping to prevent idling
	// config := newrelic.NewConfig("yeardecadeprogress", keystore.NewRelicKey)
	// app, _ := newrelic.NewApplication(config)
	// router := mux.NewRouter()
	// router.HandleFunc("/ping", getResponse).Methods("GET")
	// http.HandleFunc(newrelic.WrapHandleFunc(app, "/ping", getResponse))

	// setting port for New Relic, since they require it
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "5000"
	// }
	// http.ListenAndServe(":"+port, router)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println(<-ch)

}

func sendTweet(status string, percent int, now time.Time, client *twitter.Client) {

}

func calcYearCompleted(now time.Time, nextYear int, client *twitter.Client) {

	log.Info("Calculating year completed.......")

	difference := now.Sub(time.Date(nextYear, time.January, 1, 0, 0, 0, 0, time.UTC))

	var percent int
	if now.Year()%4 == 0 {
		percent = int(((daysInLeapYear - (difference.Hours() / -24)) / daysInLeapYear) * 100)
	}
	percent = int(((daysInYear - (difference.Hours() / -24)) / daysInYear) * 100)

	if percent > yearProgress {

		log.Info("Composting year progress teweet.......")

		yearProgress = percent
		var status string
		for i := 0; i <= (100 * .5); i = i + 4 {
			if float64(i) <= (float64(percent) * .5) {
				status += completed
			} else {
				status += notCompleted
			}
		}

		status = fmt.Sprintf("Year of %[1]d\n\n%[2]s - %[3]d%%", now.Year(), status, percent)

		log.Info("Sending year progress tweet.......")

		// Send a Tweet
		_, _, err := client.Statuses.Update(status, nil)
		if err != nil {
			log.Error(err)
		}
	}

	if yearProgress == 100 {
		yearProgress = 0
	}

}

func calcDecadeCompleted(now time.Time, decadeEnd int, client *twitter.Client) {

	log.Info("Calculating decade completed.......")

	difference := now.Sub(time.Date(decadeEnd, time.January, 1, 0, 0, 0, 0, time.UTC))

	numofLeapYears := 0
	for i := (decadeEnd - 10); i < decadeEnd; i++ {
		if i%4 == 0 {
			numofLeapYears++
		}
	}

	daysInDecade := (10-numofLeapYears)*daysInYear + (numofLeapYears * daysInLeapYear)
	percent := int(float64(((float64(daysInDecade) - difference.Hours()/-24) / float64(daysInDecade)) * 100))

	if int(percent) > decadeProgress {

		log.Info("Composting decade progress teweet.......")

		decadeProgress = percent
		var status string
		for i := 0; i <= (100 * .5); i = i + 4 {
			if float64(i) <= (float64(percent) * .5) {
				status += completed
			} else {
				status += notCompleted
			}
		}

		status = fmt.Sprintf("Decade of %[1]d\n\n%[2]s - %[3]d%%", decadeEnd-10, status, percent)

		log.Info("Sending year progress tweet.......")

		// Send a Tweet
		_, _, err := client.Statuses.Update(status, nil)
		if err != nil {
			log.Error(err)
		}
	}

	if decadeProgress == 100 {
		decadeProgress = 0
	}
}

func getResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Success")
}
