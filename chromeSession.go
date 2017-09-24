package headlessChrome

import (
	"bufio"
	"io"
	"os/exec"
)

// ChromeSession is an interactive console session with a Chrome
// instance.
type ChromeSession struct {
	input    io.Writer   // input to be written to the console
	output   io.Reader   // output coming from the console
	cliError io.Reader   // error output from the shell
	Input    chan string // incoming lines of input
	Output   chan string // outgoing lines of input
}

// WriteString writes a string to the console as if you wrote
// it and pressed enter.
func (cs *ChromeSession) writeString(s string) error {
	_, err := io.WriteString(cs.input, s)
	return err
}

// ReadAllOutput reads all output as of when read.  Each
// output is a line.
func (cs *ChromeSession) startOutputReader() {
	reader := bufio.NewScanner(cs.output)
	for reader.Scan() {
		cs.Output <- reader.Text()
	}
}

func (cs *ChromeSession) startInputForwarder() {
	for i := range cs.Input {
		cs.writeString(i)
	}
}

// Init runs things required to initalize a chrome session.
// No need to call outside of NewChromeSession (which does
// it for you)
func (cs *ChromeSession) Init() {
	cs.startOutputReader()
	cs.startInputForwarder()
}

// NewChromeSession starts a new chrome headless session.
func NewChromeSession(url string) (*ChromeSession, error) {
	var session ChromeSession
	var err error

	// /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --headless --repl --disable-gpu URL
	args := []string{
		"--headless",
		"--reply",
		"--disable-gpu",
		"--repl",
		url,
	}

	// setup the command and input/output pipes
	chromeStartString := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	cmd := exec.Command(chromeStartString, args...)
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return &session, err
	}
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return &session, err
	}
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return &session, err
	}

	// bind sessions to struct
	session.output = outPipe
	session.input = inPipe
	session.cliError = errPipe

	// make channels for communication
	session.Input = make(chan string, 500)
	session.Output = make(chan string, 500)

	// start channeling output and other requirements
	session.Init()

	// return a pointer to the new session
	return &session, err

}
