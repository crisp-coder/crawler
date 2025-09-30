package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})

	for _, page := range pages {
		writer.Write(
			[]string{
				page.URL,
				page.H1,
				page.FirstParagraph,
				strings.Join(page.OutgoingLinks, ";"),
				strings.Join(page.ImageURLs, ";"),
			})
	}

	writer.Flush()
	return nil
}
