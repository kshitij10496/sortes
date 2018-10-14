package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const URL = "http://api.aeneid.eu/sortes"

type Sortes struct {
	Book          int      `json:"book,omitempty"`
	NumberOfLines int      `json:"number_of_lines,omitempty"`
	StartLine     int      `json:"start_line.omitempty"`
	Text          []string `json:"text,omitempty"`
	Version       string   `json:"version,omitempty"`
}

func main() {
	flagVersion := flag.String("version", "dryden", `version`)
	flagLines := flag.Int("lines", 1, `number of lines to fetch`)
	flagBook := flag.Int("book", 0, `book number of Aeneid`)

	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Cannot create request:", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Cannot make request:", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("Error making the request:", res.Status)
	}

	defer res.Body.Close()

	var content Sortes
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&content)
	if err != nil {
		log.Fatal("Unable to decode the response:", err)
	}

	fmt.Println(strings.Join(content.Text, " "))
}
