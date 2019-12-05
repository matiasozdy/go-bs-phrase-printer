package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Phrase struct {
  Phrase string
}

func main() {
	var client http.Client
	var url = "https://corporatebs-generator.sameerkumar.website"
	resp, err := client.Get(url)
	if err != nil {
	    log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
	    bodyBytes, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
	        log.Fatal(err)
	    }
	    var phrase1 Phrase
	    jsonErr := json.Unmarshal(bodyBytes, &phrase1)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	    fmt.Printf("Motivational bs phrase of the day: %s\n", phrase1.Phrase)
}
}
