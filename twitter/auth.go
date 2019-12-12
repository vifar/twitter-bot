package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Keystore of the API key values
type Keystore struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
}

var keystore Keystore

// Authenticating with the twitter web client
func auth() *twitter.Client {
	getKeys()
	config := oauth1.NewConfig(keystore.ConsumerKey, keystore.ConsumerSecret)
	token := oauth1.NewToken(keystore.AccessToken, keystore.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// verify := &twitter.AccountVerifyParams{
	// 	SkipStatus:   twitter.Bool(true),
	// 	IncludeEmail: twitter.Bool(true),
	// }
	// user, _, _ := client.Accounts.VerifyCredentials(verify)
	// log.Info("User's Account:", user)

	log.Info("Authenticated & Connected.......")

	return client
}

// Getting keys from either the keys.JSON file or from environment variables
func getKeys() {
	// Reading keystore file, if not available read from environment variables
	jsonFile, err := os.Open("keys.json")
	if err != nil {
		keystore.ConsumerKey = os.Getenv("consumerKey")
		keystore.ConsumerSecret = os.Getenv("consumerSecret")
		keystore.AccessToken = os.Getenv("accessToken")
		keystore.AccessSecret = os.Getenv("accessSecret")
		return
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &keystore)
	}
	defer jsonFile.Close()

	log.Info("Retrieved Keys.......")

}
