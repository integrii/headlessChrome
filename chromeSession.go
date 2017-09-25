package headlessChrome

import (
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

	// close will cause the io workers to stop gracefully
	close(cs.session.Input)
}

// Write writes to the Session
func (cs *ChromeSession) Write(s string) {
	cs.session.Write(s)
}

// NewChromeSession starts a new chrome headless session.
func NewChromeSession(url string) (*ChromeSession, error) {
	var chromeSession *ChromeSession
	var err error

	// add url as last arg and create new session
	args := append(Args, url)
	chromeSession.session, err = interactive.NewSession(ChromeCommand, args)

	return chromeSession, err
}

// forceClose issues a force kill to the command
func (cs *ChromeSession) forceClose() {
	cs.session.ForceClose()
}
