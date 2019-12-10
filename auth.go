package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Keystore of the API key values
type Keystore struct {
	APIKey       string `json:"apiKey"`
	APISecretKey string `json:"apiSecretKey"`
}

var keystore Keystore

func auth() {

	// Reading keystore file, if not available read from environment variables
	jsonFile, err := os.Open("keys.json")
	if err != nil {
		keystore.APIKey = os.Getenv("API_KEY")
		keystore.APISecretKey = os.Getenv("API_SECRET_KEY")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &keystore)
	fmt.Println("Retrieved Keys.......")

}
