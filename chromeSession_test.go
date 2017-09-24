package headlessChrome

import (
	"fmt"
	"testing"
)

// TestHTTPScrape tests a scrape from content on httpbin.org
func TestHTTPScrape(t *testing.T) {

	Debug = false

	// make a new session
	session, err := NewChromeSession("http://google.com")
	if err != nil {
		t.Fatal(err)
	}

	session.Write(`console.log("DEBUG")\r\n`)
	session.Write(`console.log(document)`)
	session.Write(`console.log("DEBUG")`)
	session.Write(`console.log("DEBUG")`)
	session.Exit()

	// write to the session and issue an exit
	for l := range session.Output {
		t.Log("Output:", l)
		fmt.Println("Output:", l)
	}

	// b, err := ioutil.ReadAll(session.CLIError)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))
}
