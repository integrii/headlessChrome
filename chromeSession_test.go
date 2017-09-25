package headlessChrome

import (
	"fmt"
	"testing"
	"time"
)

// TestHTTPScrape tests a scrape from content on httpbin.org
func TestHTTPScrape(t *testing.T) {

	Debug = false

	// make a new session
	chromeSession, err := NewChromeSession("http://google.com")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
	chromeSession.Write(`document.documentElement.outerHTML`)
	chromeSession.Exit()

	// write to the session and issue an exit
	for l := range chromeSession.session.Output {
		t.Log("Output:", l)
		fmt.Println("Output:", l)
	}

	// b, err := ioutil.ReadAll(session.CLIError)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))
}
