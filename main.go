package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	bURL, err := url.Parse(baseURL)
	if err != nil {
		log.Println("Error parsing URL")
		os.Exit(1)
	}
	pages := make(map[string]PageData)
	var mu sync.Mutex
	var wg sync.WaitGroup
	ch := make(chan struct{}, 5)
	cfg := config{
		pages:              pages,
		baseURL:            bURL,
		mu:                 &mu,
		concurrencyControl: ch,
		wg:                 &wg,
	}

	fmt.Printf("starting crawl of: %v\n", baseURL)
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

	for key, val := range pages {
		if val.URL != "" {
			fmt.Printf("%v: %v\n", key, val)
		}
	}
}
