package headlessChrome

import (
	"fmt"

	"github.com/integrii/interactive"
)

// Debug enables debug output for this package to console
var Debug bool

// ChromeCommand is the command to execute chrome
var ChromeCommand = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

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

// ChromeSession is an interactive console session with a Chrome
// instance.
type ChromeSession struct {
	session *interactive.Session
}

// Exit exits the running command out by ossuing a 'quit'
// to the chrome console
func (cs *ChromeSession) Exit() {
	cs.session.Write(`quit`)
	cs.session.Exit()
}

// Write writes to the Session
func (cs *ChromeSession) Write(s string) {
	cs.session.Write(s)
}

// OutputPrinter prints all outputs from the output channel to the cli
func (cs *ChromeSession) OutputPrinter() {
	for l := range cs.session.Output {
		fmt.Println(l)
	}
}

// forceClose issues a force kill to the command
func (cs *ChromeSession) forceClose() {
	cs.session.ForceClose()
}

// ClickSelector calls a click() on the supplied selector
func (cs *ChromeSession) ClickSelector(s string) {
	cs.Write(`document.querySelector("` + s + `").click()`)
}

// ClickItemWithInnerHTML clicks an item that has the matching inner html
func (cs *ChromeSession) ClickItemWithInnerHTML(elementType string, s string) {
	cs.Write(`var spans = $("` + elementType + `").filter(function(idx) { return this.innerHTML.indexOf("` + s + `") == 0; }); spans[0].click()`)
}

// NewChromeSession starts a new chrome headless session.
func NewChromeSession(url string) (*ChromeSession, error) {
	var err error

	chromeSession := ChromeSession{}

	// add url as last arg and create new session
	args := append(Args, url)
	chromeSession.session, err = interactive.NewSession(ChromeCommand, args)

	return &chromeSession, err
}
