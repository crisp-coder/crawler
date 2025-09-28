package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println(err)
	}
	text := doc.Find("h1").First().Text()
	return text
}

func getFirstParagraphFromHTML(html string) string {
	htmlReader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println(err)
		return ""
	}
	text := doc.Find("p").First().Text()
	mainText := doc.Find("main").Find("p").Text()
	if mainText == "" {
		return text
	}
	return mainText
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println(err)
		return make([]string, 0), err
	}

	// Search every node for urls
	urls := make([]string, 0)
	selection := doc.Find("a")
	for _, link := range selection.Nodes {
		for _, attr := range link.Attr {
			if attr.Key == "href" && attr.Val != "" {
				urls = append(urls, attr.Val)
			}
		}
	}

	// Make urls absolute
	abs_urls := make([]string, 0)
	for _, u := range urls {
		parsed_url, err := url.Parse(u)
		if err != nil {
			log.Println(err)
			continue
		}
		resolved_url := baseURL.ResolveReference(parsed_url)
		abs_urls = append(abs_urls, resolved_url.String())
	}

	return abs_urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Println(err)
		return make([]string, 0), err
	}

	images := make([]string, 0)
	selection := doc.Find("img")
	for _, image := range selection.Nodes {
		for _, attr := range image.Attr {
			if attr.Key == "src" && attr.Val != "" {
				images = append(images, attr.Val)
			}
		}
	}

	abs_images := make([]string, 0)
	for _, u := range images {
		parsed_url, err := url.Parse(u)
		if err != nil {
			log.Println(err)
			continue
		}
		resolved_url := baseURL.ResolveReference(parsed_url)
		abs_images = append(abs_images, resolved_url.String())
	}

	return abs_images, nil
}

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	pageData := PageData{}
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		log.Println(err)
		return pageData
	}
	pageData.URL = baseURL.String()
	pageData.H1 = getH1FromHTML(html)
	pageData.FirstParagraph = getFirstParagraphFromHTML(html)

	links, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		log.Println(err)
		return pageData
	}
	pageData.OutgoingLinks = links

	images, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		log.Println(err)
		return pageData
	}
	pageData.ImageURLs = images

	return pageData
}
