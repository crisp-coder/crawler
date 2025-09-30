package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Parse args
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		log.Println(err)
		return
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		log.Println(err)
		return
	}

	bURL, err := url.Parse(baseURL)
	if err != nil {
		log.Println("Error parsing URL")
		os.Exit(1)
	}

	// Prepare config for crawler
	pages := make(map[string]PageData)
	var mu sync.Mutex
	var wg sync.WaitGroup
	ch := make(chan struct{}, maxConcurrency)
	cfg := config{
		pages:              pages,
		baseURL:            bURL,
		mu:                 &mu,
		concurrencyControl: ch,
		wg:                 &wg,
		maxPages:           maxPages,
	}

	// Launch crawler
	cfg.wg.Add(1)
	u := bURL.String()
	go func(baseURL string) {
		cfg.concurrencyControl <- struct{}{}
		defer func() {
			cfg.wg.Done()
			<-cfg.concurrencyControl
		}()
		cfg.crawlPage(baseURL)
	}(u)
	cfg.wg.Wait()

	// Print report
	for key, val := range pages {
		fmt.Printf("%v: %v\n", key, val.URL)
	}

	writeCSVReport(pages, "report.csv")
}
