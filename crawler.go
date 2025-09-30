package main

import (
	"log"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	maxPages           int
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// Stop crawling once max pages is reached
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	// Parse current URL
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Println("bad current:", err)
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
		log.Printf("Content-Type not text/html: %v\n", currentURL.String())
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

	log.Printf("starting crawl of: %v\n", rawCurrentURL)
	for _, link := range links {

		// Skip links with these prefixes
		if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
			continue
		}

		n, err := normalizeURL(link)
		if err != nil {
			continue
		}

		u, err := url.Parse(link)
		if err != nil {
			continue
		}

		// Verify next link is on same domain
		if currentURL.Hostname() != cfg.baseURL.Hostname() {
			continue
		}

		if !cfg.addPageVisit(n) {
			continue
		}

		// Recursive crawl on link
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
	return true
}
