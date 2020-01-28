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
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	log "github.com/sirupsen/logrus"
)

const notCompleted = "░"
const completed = "▓"
const daysInYear = 365
const daysInLeapYear = 366

var yearProgress = 0
var decadeProgress = 0

var sendYearTweet = false
var sendDecadeTweet = false

var decadeEnd = 0;

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
	decadeEnd = int(math.Round(float64(now.Year())/10) * 10)+10
	log.Info(decadeEnd)

	ticker := time.NewTicker(30 * time.Second)
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

	// start new relic, and setting address to ping to prevent idling
	config := newrelic.NewConfig("yeardecadeprogress", keystore.NewRelicKey)
	app, _ := newrelic.NewApplication(config)
	router := mux.NewRouter()
	router.HandleFunc("/ping", getResponse).Methods("GET")
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/ping", getResponse))

	// setting port for New Relic, since they require it
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.ListenAndServe(":"+port, router)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println(<-ch)

}

func calcYearCompleted(now time.Time, nextYear int, client *twitter.Client) {

	log.Info("Calculating year completed.......")

	difference := now.Sub(time.Date(nextYear, time.January, 1, 0, 0, 0, 0, time.UTC))

	var percent int
	if now.Year()%4 == 0 {
		percent = int(((daysInLeapYear - (difference.Hours() / -24)) / daysInLeapYear) * 100)
	}
	percent = int(((daysInYear - (difference.Hours() / -24)) / daysInYear) * 100)
	
	var status string
	if percent > yearProgress {

		log.Info("Composting year progress teweet.......")

		sendYearTweet = true
		yearProgress = percent
		for i := 0; i <= 70; i = i + 5 {
			if float64(i) <= (float64(percent) * .7) {
				status += completed
			} else {
				status += notCompleted
			}
		}

		status = fmt.Sprintf("Year of %[1]d\n\n%[2]s %[3]d%%", now.Year(), status, percent)

	} else {
		sendYearTweet = false
	}

	if yearProgress == 100 {
		yearProgress = 0
	}

	// Send a Tweet
	if sendYearTweet {
		log.Info("Sending year progress tweet.......")
		_, _, err := client.Statuses.Update(status, nil)
		if err != nil {
			log.Error(err)
		}
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

	var status string
	if percent > decadeProgress {

		log.Info("Composting decade progress teweet.......")

		sendDecadeTweet = true
		decadeProgress = percent
		for i := 0; i <= 70; i = i + 5 {
			if float64(i) <= (float64(percent) * .7) {
				status += completed
			} else {
				status += notCompleted
			}
		}

		status = fmt.Sprintf("Decade of %[1]d\n\n%[2]s %[3]d%%", decadeEnd-10, status, percent)

	} else {
		sendDecadeTweet = false
	}

	if decadeProgress == 100 {
		decadeProgress = 0
		decadeEnd += 10
	}

	// Send a Tweet
	if sendDecadeTweet {
		log.Info("Sending decade progress tweet.......")
		_, _, err := client.Statuses.Update(status, nil)
		if err != nil {
			log.Error(err)
		}
	}
}

func getResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Success")
}
