package main

import (
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
)

func cleanPath(u *url.URL) string {
	u.Path = path.Clean(u.Path)
	empty := u.Path == "" || u.Path == "/" || u.Path == "."
	if empty {
		return ""
	}
	return strings.TrimSuffix(u.Path, "/")
}

func cleanQuery(u *url.URL) string {
	clean := url.Values{}
	q := u.Query()
	for key, vals := range q {
		if strings.TrimSpace(key) == "" {
			continue
		}

		if strings.HasPrefix(strings.ToLower(key), "utm_") {
			continue
		}

		seen := map[string]struct{}{}
		for _, e := range vals {
			e = strings.TrimSpace(e)
			if e == "" {
				continue
			}
			if _, ok := seen[e]; ok {
				continue
			}
			seen[e] = struct{}{}
			clean.Add(key, e)
		}
	}

	// Sort the keys and values for consistent url output
	keys := make([]string, 0, len(clean))
	for k := range clean {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var b strings.Builder
	first := true
	for _, k := range keys {
		vals := clean[k]
		sort.Strings(vals)
		for _, v := range vals {
			if first {
				first = false
			} else {
				b.WriteByte('&')
			}
			b.WriteString(url.QueryEscape(k))
			b.WriteByte('=')
			b.WriteString(url.QueryEscape(v))
		}
	}

	return b.String()
}

func normalizeURL(u string) (string, error) {
	parsed_url, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	url_str := ""
	url_str += strings.ToLower(parsed_url.Hostname())

	if parsed_url.Port() != "80" && parsed_url.Port() != "443" && parsed_url.Port() != "" {
		url_str += ":" + parsed_url.Port()
	}

	url_str += cleanPath(parsed_url)

	query_params := parsed_url.Query()
	if len(query_params) != 0 {
		url_str += "?"
	}

	url_str += cleanQuery(parsed_url)

	fmt.Println(url_str)
	return url_str, nil
}
