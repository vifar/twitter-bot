package main

import (
	"fmt"
	"log"
)

func main() {

	fmt.Println("Retrieving Keys.......")
	client := auth()

	// Send a Tweet
	_, resp, err := client.Statuses.Update("just setting up my twttr", nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}
