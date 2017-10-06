package headlessChrome

import (
	"strings"
	"testing"
	"time"
)

// TestMainPageScrape tests a scrape from content on httpbin.org
func TestMainPageScrape(t *testing.T) {

	Debug = true

	// make a new Session
	chrome, err := NewBrowser(`google.com`)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	chrome.Write(`document.documentElement.outerHTML`)
	chrome.Exit()

	// write to the Session and issue an exit
	var googleFound bool
	for l := range chrome.Output {
		if strings.Contains(l, "google") {
			googleFound = true
		}
		t.Log(l)
	}

	if !googleFound {
		t.Fatal("Didnt find google in the output")
	}
}
