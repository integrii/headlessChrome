# headlessChrome 🤖
**MacOS only for now.**  😬

A [go](https://golang.org) package for working with headless Chrome.  Run interactive JavaScript commands on pages with go and Chrome without a GUI.  Includes a few helpful functions out of the box to query and click selector paths by their classes or content (innerHTML).

You could use this package to click buttons and scrape content on/from a website as if you were a browser, or to render pages that wouldn't be supported by other things like phantomjs or casperjs.  Especially useful for sites that use EmberJS, where the content is rendered by javascript after the HTML payload is delivered.


#### Documentation
[http://godoc.org/github.com/integrii/headlessChrome](http://godoc.org/github.com/integrii/headlessChrome)

##### Custom Flags
By default, we startup with the bare minimum flags necessary to start headless chrome and open a javascript console.  If you want more flags, like a resolution size, or a custom User-Agent, you can specify it by replacing the `Args` variable.  Just be sure to append to it so you don't kill the default flags...

```go
headlessChrome.Args = append(headlessChrome.Args,"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
headlessChrome.Args = append(headlessChrome.Args,"--window-size=1024,768")
```


#### Install
`go get github.com/integrii/headlessChrome`

#### Example

```go
// make a new session
browser, err := headlessChrome.NewBrowser(`google.com`)
if err != nil {
  t.Fatal(err)
}

// sleep while content is rendered.  You could replace this with some javascript that only returns when the content exists.
time.Sleep(time.Second * 5)

// Write all the HTML from the web site:
browser.Write(`document.documentElement.outerHTML`)

// range over all the output that comes from the browser with the string reader package
for len(browser.Session.Output) > 0 {
  fmt.Println(<-browser.Session.Output)
}

// click some span element from the page by its text content
browser.ClickItemWithInnerHTML("span", "Google Search",0)

// drain all the output so we can pay attention to the next line
for len(c) > 0 {
  <-c
}

// select the content of something by its css classes
browser.GetContentOfItemWithClasses("button arrow bold",0)

// read the selected stuff from the console by picking the next item from the channel
consoleOut := <- browser.Session.Output
fmt.Println(consoleOut)

// End the session by writing quit
browser.Exit()
```


#### Contributing

Please send pull requests!  It would be good to have support for more OSes or more handy helpers to run more commonly used javascript code easily.
