package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Keystore of the API key values
type Keystore struct {
	APIKey       string `json:"apiKey"`
	APISecretKey string `json:"apiSecretKey"`
}

var keystore Keystore

// Authenticating with the twitter web client
func auth() *twitter.Client {
	getKeys()
	config := oauth1.NewConfig("consumerKey", keystore.APIKey)
	token := oauth1.NewToken("accessToken", keystore.APISecretKey)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	fmt.Println("Authenticated & Connected.......")

	return client
}

// Getting keys from either the keys.JSON file or from environment variables
func getKeys() {
	// Reading keystore file, if not available read from environment variables
	jsonFile, err := os.Open("keys.json")
	if err != nil {
		keystore.APIKey = os.Getenv("API_KEY")
		keystore.APISecretKey = os.Getenv("API_SECRET_KEY")
		return
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &keystore)
	}
	defer jsonFile.Close()

	fmt.Println("Retrieved Keys.......")

}
