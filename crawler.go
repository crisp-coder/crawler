package main

import (
	"log"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// Parse current URL
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Println("bad current:", err)
		return
	}

	// Verify next link is on same domain
	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		log.Println("Url is on different domain:")
		return
	}

	// Normalize the current url for page caching
	norm, err := normalizeURL(currentURL.String())
	if err != nil {
		log.Println("normalize:", err)
		return
	}

	// Fetch the resource
	html, ct, err := getHTMLWithType(currentURL.String())
	if err != nil {
		log.Println("fetch:", err)
		return
	}

	// Verify the requested resource is html
	if !strings.HasPrefix(ct, "text/html") {
		return
	}

	// Update the data in the pages map
	pageData := extractPageData(html, norm)
	cfg.mu.Lock()
	cfg.pages[norm] = pageData
	cfg.mu.Unlock()

	// Find all links in the html
	links, err := getURLsFromHTML(html, cfg.baseURL)
	if err != nil {
		log.Println("parse links:", err)
		return
	}

	// For each link, recurse if the link is not already in the pages map
	for _, link := range links {

		if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
			log.Printf("skipping link: %v\n", link)
			continue
		}

		n, err := normalizeURL(link)
		if err != nil {
			log.Printf("normalize: %v\n", err)
			continue
		}

		if !cfg.addPageVisit(n) {
			log.Printf("skipping link: %v\n", n)
			continue
		}

		u, err := url.Parse(link)
		if err != nil {
			log.Printf("url parse: %v\n", err)
			continue
		}

		// Launch go routine
		cfg.wg.Add(1)
		go func(urlLink string) {
			cfg.concurrencyControl <- struct{}{}
			defer func() {
				cfg.wg.Done()
				<-cfg.concurrencyControl
			}()
			cfg.crawlPage(urlLink)
		}(u.String())
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	}
	cfg.pages[normalizedURL] = PageData{}
	return true
}
