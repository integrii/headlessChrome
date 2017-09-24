package headlessChrome

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
)

// Debug enables debug output for this package to console
var Debug bool

const expectedFirstLine = "Type a Javascript expression to evaluate or \"quit\" to exit."

// ChromeSession is an interactive console session with a Chrome
// instance.
type ChromeSession struct {
	stdIn  io.Writer   // input to be written to the console
	stdOut io.Reader   // output coming from the console
	stdErr io.Reader   // error output from the shell
	Input  chan string // incoming lines of input
	Output chan string // outgoing lines of input
	cmd    *exec.Cmd   // cmd that holds this chrome instance
	pty    *os.File    // the tty for the session
}

// WriteString writes a string to the console as if you wrote
// it and pressed enter.
func (cs *ChromeSession) writeString(s string) error {
	if Debug {
		fmt.Println("Writing string:", s)
	}
	_, err := io.WriteString(cs.stdIn, s+"\r\n")
	// _, err := cs.pty.WriteString(s + "\r\n")
	return err
}

// ReadAllOutput reads all output as of when read.  Each
// output is a line.
func (cs *ChromeSession) startOutputReader() {
	reader := bufio.NewScanner(cs.pty)
	if Debug {
		fmt.Println("Output reader looking for output")
	}
	for reader.Scan() {
		if Debug {
			fmt.Println("Output reader got text:", reader.Text())
		}
		cs.Output <- reader.Text()
		if Debug {
			fmt.Println("Reader passed text to outut channel:", reader.Text())
		}
	}
}

// startErrorReader starts an error reader that outputs
// to the output channel
func (cs *ChromeSession) startErrorReader() {
	reader := bufio.NewScanner(cs.stdErr)
	if Debug {
		fmt.Println("Error reader looking for output")
	}
	for reader.Scan() {
		if Debug {
			fmt.Println("Error reader got text:", reader.Text())
		}
		cs.Output <- reader.Text()
		if Debug {
			fmt.Println("Error reader passed output to channel:", reader.Text())
		}
	}
}

// startInputForwarder starts a forwarder of input channel
// to running session
func (cs *ChromeSession) startInputForwarder() {
	for l := range cs.Input {
		if Debug {
			fmt.Println("Got request to write string:", l)
		}
		cs.writeString(l)
	}
}

// Exit exits the running command out by ossuing a 'quit'
// to the chrome console
func (cs *ChromeSession) Exit() {
	cs.Write(`quit`)

	// close will cause the io workers to stop gracefully
	close(cs.Input)
}

// Init runs things required to initalize a chrome session.
// No need to call outside of NewChromeSession (which does
// it for you)
func (cs *ChromeSession) Init() {
	go cs.startOutputReader()
	go cs.startErrorReader()
	go cs.startInputForwarder()
	go cs.closeWhenCompleted()
}

// Write writes an output line into the session
func (cs *ChromeSession) Write(s string) {
	cs.Input <- s
}

// closeWhenCompleted closes ouput channels to cause readers to
// end gracefully when the command completes
func (cs *ChromeSession) closeWhenCompleted() {

	if Debug {
		fmt.Println("Spawned chrome as PID", cs.cmd.Process.Pid)
	}

	cs.cmd.Wait()
	if Debug {
		fmt.Println("Command exited. Closing channels.")
	}
	close(cs.Output)

	cs.forceClose() // when complete, make sure the PID dies (chrome never does on its own as of writing)
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
	session.cmd = exec.Command(chromeStartString, args...)
	errPipe, err := session.cmd.StderrPipe()
	if err != nil {
		return &session, err
	}
	inPipe, err := session.cmd.StdinPipe()
	if err != nil {
		return &session, err
	}
	outPipe, err := session.cmd.StdoutPipe()
	if err != nil {
		return &session, err
	}

	// bind sessions to struct
	session.stdOut = outPipe
	session.stdIn = inPipe
	session.stdErr = errPipe

	// make channels for input and outut communication to the process
	session.Input = make(chan string, 1)
	session.Output = make(chan string, 5000)

	if Debug {
		fmt.Println("Starting command:", session.cmd.Args)
	}

	// kick off the command and ensure it closes when done
	session.pty, _ = pty.Start(session.cmd)
	if err != nil {
		return &session, err
	}

	// start channeling output and other requirements
	session.Init()

	// read the first item off the session (the welcome banner)
	// before returning
	firstLine := <-session.Output
	if !strings.Contains(firstLine, expectedFirstLine) {
		return &session, errors.New("Was unable to fetch proper console startup line.  Got: " + firstLine + " but expected " + expectedFirstLine)
	}

	// command is online and healthy, return to the user
	return &session, err

}

// forceClose issues a force kill to the command
func (cs *ChromeSession) forceClose() {
	cs.cmd.Process.Kill()
}
