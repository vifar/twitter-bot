package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	log.Info("Retrieving Keys.......")
	client := auth()

	// Send a Tweet
	_, resp, err := client.Statuses.Update("just setting up my twttr", nil)
	if err != nil {
		log.Warn(err)
	}
	fmt.Println(resp)
}
