package log

import (
	"io/ioutil"
	logStd "log"
	"os"
	"sync"
)

const (
	DebugLevel = iota
	InfoLevel
	ErrorLevel
	// https://golang.org/src/log/log.go?s=8938:8966#L35
	Ldate         = logStd.Ldate                // the date in the local time zone: 2009/01/23
	Ltime         = logStd.Ltime                // the time in the local time zone: 01:23:23
	Lmicroseconds = logStd.Lmicroseconds        // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile     = logStd.Llongfile            // full file name and line number: /a/b/c/d.go:23
	Lshortfile    = logStd.Lshortfile           // final file name element and line number: d.go:23. overrides Llongfile
	LUTC          = logStd.LUTC                 // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = logStd.Ldate | logStd.Ltime // initial values for the standard logger
)

var (
	Debug *logStd.Logger
	Info  *logStd.Logger
	Error *logStd.Logger
	mu    *sync.Mutex = &sync.Mutex{}
	// outputs is initializes the writers to the InfoLevel and sets
	// the io writers to ioutil.Discard, os.Stdout, os.Stderr
	initOutputs = func() {
		Debug.SetOutput(ioutil.Discard)
		Info.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
	}
)

func init() {
	mu.Lock()
	defer mu.Unlock()

	Debug = logStd.New(ioutil.Discard, "DEBU: ", Ldate|Ltime|Lshortfile)
	Info = logStd.New(os.Stdout, "INFO: ", Ldate|Ltime|Lshortfile)
	Error = logStd.New(os.Stderr, "ERRO: ", Ldate|Ltime|Lshortfile)
}

// SetOutputs uses an anonomous function to initialize the writers for Debug, Info, Error
func SetOutputs(fn func()) {
	mu.Lock()
	defer mu.Unlock()
	initOutputs = fn
	initOutputs()
}

// SetFlags initializes the flags for Debug, Info, Error
func SetFlags(flags int) {
	Debug.SetFlags(flags)
	Info.SetFlags(flags)
	Error.SetFlags(flags)
}

// SetPrefix initializes the prefix for Debug, Info, Error
// It will always put the level first
func SetPrefix(prefix string) {
	Debug.SetPrefix("DEBU " + prefix + ":")
	Info.SetPrefix("INFO " + prefix + ":")
	Error.SetPrefix("ERRO " + prefix + ":")
}

// SetLevel initializes writers below the level to io.Discard.
func SetLevel(level int) {
	initOutputs()

	switch level {
	case InfoLevel:
		Debug.SetOutput(ioutil.Discard)
	case ErrorLevel:
		Debug.SetOutput(ioutil.Discard)
		Info.SetOutput(ioutil.Discard)
	}
}
