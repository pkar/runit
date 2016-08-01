# log

That's it, just simple leveled logging with the Go std lib log.

Usage:

```go
package main

import (
	"io"
	"os"

	"github.com/pkar/log"
)

func main() {
	logFile, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error.Fatalln("Failed to open log file", err)
	}
	log.SetFlags(log.Lshortfile)

	// default outputs
	log.Debug.Println("default debug")
	log.Info.Println("default info")
	log.Error.Printf("%s", "default error")

	// Create a default output setter optionally. The initial one uses
	// ioutil.Discard, os.Stdout, os.Stderr
	var newOuts = func() {
		log.Debug.SetOutput(os.Stdout)
		log.Info.SetOutput(io.MultiWriter(os.Stdout, logFile))
		log.Error.SetOutput(os.Stderr)
	}
	log.SetOutputs(newOuts)
	// Set the level dynamically
	log.SetLevel(log.ErrorLevel)
	log.Debug.Println("debug nope")
	log.Info.Println("info nope")
	log.Error.Println("error yes")

	// Reset to debug and change flags
	log.SetFlags(log.Llongfile)
	log.SetLevel(log.DebugLevel)
	log.Debug.Println("debug yes")
	log.Info.Println("info yes")
	log.Error.Fatal("fatal")
}
```

```
$ go run example/main.go
INFO: 2016/02/13 08:24:00 main.go:18: default info
ERRO: 2016/02/13 08:24:00 main.go:19: default error
ERRO: main.go:34: error yes
DEBU: /Volumes/Data/dropbox/development/euler/src/github.com/pkar/log/example/main.go:39: debug yes
INFO: /Volumes/Data/dropbox/development/euler/src/github.com/pkar/log/example/main.go:40: info yes
ERRO: /Volumes/Data/dropbox/development/euler/src/github.com/pkar/log/example/main.go:41: fatal
exit status 1
```
