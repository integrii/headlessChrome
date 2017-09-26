package headlessChrome

import (
	"fmt"
	"strconv"

	"github.com/integrii/interactive"
)

// Debug enables debug output for this package to console
var Debug bool

// ChromePath is the command to execute chrome
var ChromePath = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

// Args are the args that will be used to start chrome
var Args = []string{
	"--headless",
	"--disable-gpu",
	"--repl",
	// "--dump-dom",
	// "--window-size=1024,768",
	// "--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	// "--verbose",
}

const expectedFirstLine = "Type a Javascript expression to evaluate or \"quit\" to exit."

// ChromeSession is an interactive console Session with a Chrome
// instance.
type ChromeSession struct {
	Session *interactive.Session
}

// Exit exits the running command out by ossuing a 'quit'
// to the chrome console
func (cs *ChromeSession) Exit() {
	cs.Session.Write(`quit`)
	cs.Session.Exit()
}

// Write writes to the Session
func (cs *ChromeSession) Write(s string) {
	cs.Session.Write(s)
}

// OutputPrinter prints all outputs from the output channel to the cli
func (cs *ChromeSession) OutputPrinter() {
	for l := range cs.Session.Output {
		fmt.Println(l)
	}
}

// forceClose issues a force kill to the command
func (cs *ChromeSession) forceClose() {
	cs.Session.ForceClose()
}

// ClickSelector calls a click() on the supplied selector
func (cs *ChromeSession) ClickSelector(s string) {
	cs.Write(`document.querySelector("` + s + `").click()`)
}

// ClickItemWithInnerHTML clicks an item that has the matching inner html
func (cs *ChromeSession) ClickItemWithInnerHTML(elementType string, s string, itemIndex int) {
	cs.Write(`var spans = $("` + elementType + `").filter(function(idx) { return this.innerHTML.indexOf("` + s + `") == 0; }); spans[` + strconv.Itoa(itemIndex) + `].click()`)
}

// GetContentOfItemWithClasses fetches the content of the element with the specified classes
func (cs *ChromeSession) GetContentOfItemWithClasses(classes string, itemIndex int) {
	cs.Write(`var x = document.getElementsByClassName("` + classes + `");x[` + strconv.Itoa(itemIndex) + `].innerHTML`)
}

// ClickItemWithClasses clicks on the first item it finds with the provided classes.
// Multiple classes are separated by spaces
func (cs *ChromeSession) ClickItemWithClasses(classes string, itemIndex int) {
	cs.Write(`var x = document.getElementsByClassName("` + classes + `");x[` + strconv.Itoa(itemIndex) + `].click()`)
}

// SetTextByID sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByID(divID string, itemIndex int, text string) {
	cs.Write(`var x = document.getElementsById("` + divID + `");x[` + strconv.Itoa(itemIndex) + `].innerHTML = "` + text + `"`)
}

// SetTextByClasses sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`var x = document.getElementsByClassName("` + classes + `");x[` + strconv.Itoa(itemIndex) + `].innerHTML = "` + text + `"`)
}

// SetInputTextByClasses sets the input text for an input field
func (cs *ChromeSession) SetInputTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`var x = document.getElementsByClassName("` + classes + `");x[` + strconv.Itoa(itemIndex) + `].value = "` + text + `"`)
}

// NewBrowser starts a new chrome headless Session.
func NewBrowser(url string) (*ChromeSession, error) {
	var err error

	chromeSession := ChromeSession{}

	// add url as last arg and create new Session
	args := append(Args, url)
	chromeSession.Session, err = interactive.NewSession(ChromePath, args)

	return &chromeSession, err
}
