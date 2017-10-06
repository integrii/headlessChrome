package main

import (
	"fmt"
	"time"

	"github.com/integrii/headlessChrome"
)

func main() {

	// set headless chrome package into debug mode (prints to stdout)
	// headlessChrome.Debug = true
	// interactive.Debug = true

	// set the path to the docker container chrome executable
	// headlessChrome.ChromePath = headlessChrome.ChromePathDocker
	headlessChrome.ChromePath = `/opt/google/chrome-unstable/chrome`

	// set some additional arguments for when starting chrome
	headlessChrome.Args = append(headlessChrome.Args, "--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	headlessChrome.Args = append(headlessChrome.Args, "--window-size=1024,768")
	headlessChrome.Args = append(headlessChrome.Args, "--no-sandbox")

	// make a new session
	browser, err := headlessChrome.NewBrowser(`http://ebay.com`)
	if err != nil {
		panic(err)
	}
	// Close the browser process when this func returns
	defer browser.Exit()

	// sleep while content is rendered.  You could replace this
	// with some javascript that only returns when the
	// content exists to be safer.
	time.Sleep(time.Second * 5)

	// Query all the HTML from the web site for fun
	browser.Write("document.documentElement.outerHTML")
	time.Sleep(time.Second * 1)

	// loop over all the output that came from the ouput channel
	// and print it to the console
	for len(browser.Output) > 0 {
		fmt.Println(<-browser.Output)
	}

}
