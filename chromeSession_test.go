package headlessChrome

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestMainPageScrape tests a scrape from content on httpbin.org
func TestMainPageScrape(t *testing.T) {

	Debug = false

	// make a new session
	chrome, err := NewBrowser(`google.com`)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	chrome.Write(`document.documentElement.outerHTML`)
	chrome.Exit()

	// write to the session and issue an exit
	var googleFound bool
	for l := range chrome.session.Output {
		if strings.Contains(l, "google") {
			googleFound = true
		}
		fmt.Println(l)
	}

	// b, err := ioutil.ReadAll(session.CLIError)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))

	if !googleFound {
		t.Fatal("Didnt find google in the output")
	}
}
