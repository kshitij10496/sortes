package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const URL = "http://api.aeneid.eu/sortes"

type Sortes struct {
	Book          int      `json:"book,omitempty"`
	NumberOfLines int      `json:"number_of_lines,omitempty"`
	StartLine     int      `json:"start_line,omitempty"`
	Text          []string `json:"text,omitempty"`
	Version       string   `json:"version,omitempty"`
}

var (
	flagVersion = flag.String("version", "dryden", `version`)
	flagLines   = flag.Int("lines", 5, `number of lines to fetch`)
	flagBook    = flag.Int("book", 0, `book number of Aeneid`)
)

func main() {
	flag.Parse()

	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Cannot create request:", err)
	}

	q := req.URL.Query()
	q.Add("version", *flagVersion)
	q.Add("lines", strconv.Itoa(*flagLines))
	if *flagBook != 0 {
		q.Add("book", strconv.Itoa(*flagBook))
	}

	req.URL.RawQuery = q.Encode()

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

	if content.NumberOfLines != *flagLines {
		log.Fatal("Error from the server")
	}

	fmt.Println()
	fmt.Printf("%25s %v\n", "BOOK", content.Book)
	fmt.Println()
	for i, text := range content.Text {
		fmt.Printf("%4v %s\n", content.StartLine+i, text)
	}
	fmt.Println()
	fmt.Printf("%50s: %s\n", "Aeneid", "Virgil")
	fmt.Printf("%50s: %s\n", "version", content.Version)
}
