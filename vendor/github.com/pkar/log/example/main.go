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

	// default outputs
	log.Debug.Println("default debug")
	log.Info.Println("default info")
	log.Error.Printf("%s", "default error")

	// set prefix on debug
	log.SetPrefix("sometext")
	log.Info.Println("info prefix")
	log.Error.Println("error prefix")
	// reset it
	log.SetPrefix("")

	// Create a default output setter optionally. The initial one uses
	// ioutil.Discard, os.Stdout, os.Stderr
	var newOuts = func() {
		log.Debug.SetOutput(os.Stdout)
		log.Info.SetOutput(io.MultiWriter(os.Stdout, logFile))
		log.Error.SetOutput(os.Stderr)
	}
	log.SetFlags(log.Lshortfile)
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
