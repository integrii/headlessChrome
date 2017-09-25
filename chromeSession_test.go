package headlessChrome

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFetchAvailability(t *testing.T) {
	// make a new session
	chromeSession, err := NewChromeSession(`https://secure.rocket-rez.com/RocketWeb/?eid=8d687efc97053cc6`)
	if err != nil {
		t.Fatal(err)
	}

	go chromeSession.OutputPrinter()

	// load the main react page
	time.Sleep(time.Second * 5)
	// chromeSession.Write(`document.documentElement.outerHTML`)
	// choose bremerton to seattle one way by the text content
	chromeSession.ClickItemWithInnerHTML("span", "One Way Trip: Seattle to Bremerton")
	time.Sleep(time.Second * 1)

	// click add one adult plus button
	chromeSession.ClickItemWithClasses("plus icon")
	time.Sleep(time.Second * 1)

	// click choose times
	chromeSession.ClickItemWithClasses("ui right floated huge primary button")
	time.Sleep(time.Second * 1)

	// click the choose times calendar dropdown
	chromeSession.ClickItemWithClasses("ember-view ember-text-field datepicker picker__input")
	time.Sleep(time.Second * 1)

	// fetch the month on the calendar
	chromeSession.GetContentOfItemWithClasse("picker__month")
	time.Sleep(time.Second * 1)

	// exit
	chromeSession.Exit()
}

// TestMainPageScrape tests a scrape from content on httpbin.org
func TestMainPageScrape(t *testing.T) {

	Debug = false

	// make a new session
	chromeSession, err := NewChromeSession(`https://secure.rocket-rez.com/RocketWeb/?eid=8d687efc97053cc6`)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	chromeSession.Write(`document.documentElement.outerHTML`)
	chromeSession.Exit()

	// write to the session and issue an exit
	var bremertonFound bool
	for l := range chromeSession.session.Output {
		if strings.Contains(l, "Bremerton") {
			bremertonFound = true
		}
		fmt.Println(l)
	}

	// b, err := ioutil.ReadAll(session.CLIError)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(b))

	if !bremertonFound {
		t.Fatal("Didnt find Bremerton in the output")
	}
}
