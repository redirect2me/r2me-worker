package main

import (
	"net"
	"net/http"
	"sort"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

var recentCache *lru.Cache[string, string]

func getHostWithoutPort(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		// If SplitHostPort fails, r.Host likely had no port
		return r.Host
	}
	return host
}

func RecentAddHelper(r *http.Request, result *MapResult) {
	host := getHostWithoutPort(r)
	if host == "" {
		host = "unknown_host"
	} else {
		host = strings.ToLower(host)
		host = strings.TrimSuffix(host, ".")
		// get etld+1 for the host
		domain, domainErr := publicsuffix.EffectiveTLDPlusOne(host)
		if domainErr != nil {
			Logger.Warn("Failed to get eTLD+1 for host", "host", host, "error", domainErr)
		} else {
			host = domain
		}
		// convert from punycode to unicode
		unicodeHost, unicodeErr := idna.ToUnicode(host)
		if unicodeErr != nil {
			Logger.Warn("Failed to convert host to unicode", "host", host, "error", unicodeErr)
		} else {
			host = unicodeHost
		}
	}

	resultCode := "null"
	if result != nil {
		resultCode = result.ResultCode
	}

	RecentAdd(host, resultCode)
}

func RecentAdd(host, result string) {

	if recentCache == nil {
		recentCache, _ = lru.New[string, string](100)
	}
	existing, ok := recentCache.Get(host)
	if ok {
		if existing != result {
			Logger.Info("Updating recent cache for host", "host", host, "old_result", existing, "new_result", result)
			recentCache.Add(host, result)
		}
	} else {
		recentCache.Add(host, result)
	}
}

type RecentData struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Size    int      `json:"size"`
	Data    []string `json:"data"`
}

func RecentHandler(w http.ResponseWriter, r *http.Request) {
	result := RecentData{}

	result.Success = true
	result.Message = "OK"
	if recentCache != nil {
		for _, key := range recentCache.Keys() {
			value, _ := recentCache.Get(key)
			result.Data = append(result.Data, key+"="+value)
		}
	}
	sort.Strings(result.Data)
	result.Size = len(result.Data)

	HandleJson(w, r, result)
}
