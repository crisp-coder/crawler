package main

import (
	"fmt"
	"io"
	"net/http"
)

func getHTMLWithType(u string) (string, string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", resp.Header.Get("Content-Type"), fmt.Errorf("status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.Header.Get("Content-Type"), err
	}
	return string(b), resp.Header.Get("Content-Type"), nil
}
