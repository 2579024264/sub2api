package service

import (
	"hash/fnv"
	"net/http"
	"strings"
)

const (
	defaultDesktopUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	usAcceptLanguage        = "en-US,en;q=0.9"
)

var desktopUserAgentPool = []string{
	defaultDesktopUserAgent,
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
}

// AccountUserAgent returns the browser UA selected by the account fingerprint group.
// Empty fingerprint groups intentionally use a default UA as a stable fallback.
func AccountUserAgent(account *Account) string {
	group := ""
	if account != nil && account.FingerprintGroupMark != nil {
		group = strings.TrimSpace(*account.FingerprintGroupMark)
	}
	if group == "" {
		return defaultDesktopUserAgent
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.ToLower(group)))
	return desktopUserAgentPool[int(h.Sum32())%len(desktopUserAgentPool)]
}

// ApplyAccountRequestHeaders enforces the browser fingerprint UA and US language
// headers for account-originated upstream requests.
func ApplyAccountRequestHeaders(req *http.Request, account *Account) {
	if req == nil {
		return
	}
	ApplyAccountHeaderValues(req.Header, account)
}

func ApplyAccountHeaderValues(header http.Header, account *Account) {
	if header == nil {
		return
	}
	header.Set("User-Agent", AccountUserAgent(account))
	header.Set("Accept-Language", usAcceptLanguage)
	header.Set("Content-Language", "en-US")
	header.Set("Sec-CH-UA-Lang", `"en-US"`)
	header.Set("X-Language", "en-US")
}
