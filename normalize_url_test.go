package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		// scheme removed
		{"remove scheme https", "https://blog.boot.dev/path", "blog.boot.dev/path"},
		{"remove scheme http", "http://blog.boot.dev/path", "blog.boot.dev/path"},

		// host case
		{"lowercase host", "https://BLOG.BOOT.DEV/Path", "blog.boot.dev/Path"},

		// trailing slash
		{"drop trailing slash", "https://blog.boot.dev/path/", "blog.boot.dev/path"},
		{"keep root slash", "https://blog.boot.dev/", "blog.boot.dev"},

		// default ports
		{"drop default http port", "http://blog.boot.dev:80/path", "blog.boot.dev/path"},
		{"drop default https port", "https://blog.boot.dev:443/path", "blog.boot.dev/path"},
		{"keep non-default port", "https://blog.boot.dev:8443/path", "blog.boot.dev:8443/path"},

		// user info
		{"drop user info", "https://user:pass@blog.boot.dev/path", "blog.boot.dev/path"},

		// fragments
		{"drop fragment", "https://blog.boot.dev/path#section", "blog.boot.dev/path"},

		// queries: order-insensitive
		{"sort query params", "https://blog.boot.dev/path?b=2&a=1", "blog.boot.dev/path?a=1&b=2"},
		{"dedupe query keys", "https://blog.boot.dev/path?a=1&a=1&b=2", "blog.boot.dev/path?a=1&b=2"},
		{"drop empty params", "https://blog.boot.dev/path?a=&b=2", "blog.boot.dev/path?b=2"},
		{"keep meaningful params", "https://blog.boot.dev/article?id=123", "blog.boot.dev/article?id=123"},
		{"drop tracking params", "https://blog.boot.dev/path?utm_source=x&b=2", "blog.boot.dev/path?b=2"},

		// oddities
		{"normalize dot segments", "https://blog.boot.dev/a/./b/../c/", "blog.boot.dev/a/c"},
		{"no path becomes empty", "https://blog.boot.dev", "blog.boot.dev"},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Fatalf("Test %d - '%s' unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %d - %s expected: %q, got: %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}
