package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// RandomWord is the reutrn obect for the getRandomWord api
type RandomWord struct {
	Word    string    `json:"word"`
	Results []Results `json:"results"`
}

// Results is an array within RandomWord
type Results struct {
	Definition   string   `json:"definition"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Synonyms     []string `json:"synonyms"`
	TypeOf       []string `json:"typeOf"`
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

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&randomWord)
	if err != nil {
		log.Errorln(nil)
	}

	log.Info("Received Random Word -->", randomWord.Word, randomWord.Results)

}
