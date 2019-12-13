package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// RandomWord is the reutrn obect for the getRandomWord api
type RandomWord struct {
	word    string
	results []Results
}

// Results is an array within RandomWord
type Results struct {
	definition   string
	partOfSpeech string
	synonyms     []string
	typeOf       []string
}

var randomWord RandomWord

func getRandomWord() {
	url := "https://wordsapiv1.p.rapidapi.com/words/?random=true"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorln(err)
	}

	req.Header.Add("x-rapidapi-host", "wordsapiv1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", keystore.WordsAPIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(res)
	}
	defer res.Body.Close()

	log.Info(res.Body)
	byteValue, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(byteValue, &randomWord)
	if err != nil {
		log.Errorln(nil)
	}

	log.Info("Received Random Word -->", randomWord)

}
