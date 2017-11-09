package headlessChrome

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/integrii/interactive"
)

// Debug enables debug output for this package to console
var Debug bool

// BrowserStartupTime is how long chrome has to startup the console
// before we consider it a failure
var BrowserStartupTime = time.Second * 20

// ChromePath is the command to execute chrome
var ChromePath = ChromePathMacOS

// ChromePathMacOS is where chrome normally lives on MacOS
var ChromePathMacOS = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

// ChromePathDocker is where chrome normally lives in the project's docker container
var ChromePathDocker = `/opt/google/chrome-unstable/chrome`

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

const expectedFirstLine = `Type a Javascript expression to evaluate or "quit" to exit.`
const promptPrefix = `>>>`

// outputSanitizer puts output coming from the consolw that
// does not begin with the input prompt into the session
// output channel
func (cs *ChromeSession) outputSanitizer() {
	for text := range cs.Session.Output {
		debug("raw output:", text)
		if !strings.HasPrefix(text, promptPrefix) {
			cs.Output <- text
		}
	}
}

// ChromeSession is an interactive console Session with a Chrome
// instance.
type ChromeSession struct {
	Session *interactive.Session
	Output  chan string
	Input   chan string
}

// Exit exits the running command out by ossuing a 'quit'
// to the chrome console
func (cs *ChromeSession) Exit() {
	cs.Session.Write(`;quit`)
	cs.Session.Exit()  // exit the process with an interrupt signal
	cs.Session.Close() // close the tty session
}

// Write writes to the Session
func (cs *ChromeSession) Write(s string) {
	debug("write:", s)
	cs.Session.Write(s)
}

// outputPrinter prints all outputs from the output channel to the cli
func (cs *ChromeSession) outputPrinter() {
	for l := range cs.Session.Output {
		debug("read:", l)
		fmt.Println(l)
	}
}

// ForceClose issues a force kill to the command
func (cs *ChromeSession) ForceClose() {
	cs.Session.ForceClose()
}

// ClickSelector calls a click() on the supplied selector
func (cs *ChromeSession) ClickSelector(s string) {
	cs.Write(`document.querySelector("` + s + `").click()`)
}

// ClickItemWithInnerHTML clicks an item that has the matching inner html
func (cs *ChromeSession) ClickItemWithInnerHTML(elementType string, s string, itemIndex int) {
	cs.Write(`var x = $("` + elementType + `").filter(function(idx) { return this.innerHTML == "` + s + `"});x[` + strconv.Itoa(itemIndex) + `].click()`)
}

// GetItemWithInnerHTML fetches the item with the specified innerHTML content
func (cs *ChromeSession) GetItemWithInnerHTML(elementType string, s string, itemIndex int) {
	cs.Write(`var x = $("` + elementType + `").filter(function(idx) { return this.innerHTML == "` + s + `"});x[` + strconv.Itoa(itemIndex) + `]`)
}

// GetContentOfItemWithClasses fetches the content of the element with the specified classes
func (cs *ChromeSession) GetContentOfItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].innerHTML`)
}

// GetValueOfItemWithClasses returns the form value of the specified item
func (cs *ChromeSession) GetValueOfItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].value`)
}

// GetContentOfItemWithSelector gets the content of an element with the specified selector
func (cs *ChromeSession) GetContentOfItemWithSelector(selector string) {
	cs.Write(`document.querySelector("` + selector + `").innerHTML()`)
}

// ClickItemWithClasses clicks on the first item it finds with the provided classes.
// Multiple classes are separated by spaces
func (cs *ChromeSession) ClickItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].click()`)
}

// SetTextByID sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByID(id string, text string) {
	cs.Write(`document.getElementById("` + id + `").innerHTML = "` + text + `"`)
}

// ClickItemWithID clicks an item with the specified id
func (cs *ChromeSession) ClickItemWithID(id string) {
	cs.Write(`document.getElementById("` + id + `").click()`)
}

// SetTextByClasses sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].innerHTML = "` + text + `"`)
}

// SetInputTextByClasses sets the input text for an input field
func (cs *ChromeSession) SetInputTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].value = "` + text + `"`)
}

// NewBrowserWithTimeout starts a new chrome headless session
// but limits how long it can run before its killed forcefully.
// A time limit of 0 means there is not a time limit
func NewBrowserWithTimeout(url string, timeout time.Duration) (*ChromeSession, error) {
	var err error

	debug("Creating a new browser pointed to", url)

	chromeSession := ChromeSession{}
	chromeSession.Output = make(chan string, 5000)

	// add url as last arg and create new Session
	args := append(Args, url)
	debug(ChromePath, args)
	chromeSession.Session, err = interactive.NewSessionWithTimeout(ChromePath, args, timeout)
	if err != nil {
		return &chromeSession, err
	}

	// map output and input channels for easy use
	chromeSession.Input = chromeSession.Session.Input
	go chromeSession.outputSanitizer()

	// wait for the console ready line from the browser
	// and if it does not start in time, throw an error
	startupTime := time.NewTimer(BrowserStartupTime)
	for {
		select {
		case <-startupTime.C:
			debug("ERROR: Browser failed to start before browser startup time cutoff")
			chromeSession.ForceClose() // force cloe the session because it failed
			err = errors.New("Chrome console failed to init in the alotted time")
			return &chromeSession, err
		case line := <-chromeSession.Output:
			if strings.Contains(line, expectedFirstLine) {
				debug("Chrome console REPL ready")
				return &chromeSession, err
			}
			debug("WARNING: Unespected first line when initializing headless Chrome console:", line)
		}
	}
}

// NewBrowser starts a new chrome headless Session.
func NewBrowser(url string) (*ChromeSession, error) {
	return NewBrowserWithTimeout(url, 0)
}

func debug(s ...interface{}) {
	if Debug {
		fmt.Println(s...)
	}
}
