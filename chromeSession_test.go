package headlessChrome

import (
	"fmt"
	"testing"
	"time"
)

// TestHTTPScrape tests a scrape from content on httpbin.org
func TestHTTPScrape(t *testing.T) {
	session, err := NewChromeSession("http://google.com")
	if err != nil {
		t.Fatal(err)
	}

	// write to the session
	session.Input <- "console.log(\"DEBUG\")"

	time.Sleep(time.Second * 2)

	t.Log(len(session.Output))

	for l := range session.Output {
		fmt.Println(l)
	}

	// b, err := ioutil.ReadAll(session.CLIError)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))
}
