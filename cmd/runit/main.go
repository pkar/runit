package main

import (
	"flag"
	"log"
	"os"

	"github.com/pkar/runit"
)

func main() {
	cmd := flag.String("cmd", "", "command to run *required")
	alive := flag.Bool("alive", false, `try to keep the command alive if it dies, you would use this for long running services like a server *optional`)
	watchPath := flag.String("watch", "", "path to directory or file to watch and restart cmd, the command will be run on startup unless wait is specified *optional")
	wait := flag.Bool("wait", false, "used with watch, this will wait for file changes and then run the cmd given *optional")
	loglevel := flag.Int("loglevel", 1, "logging level 1 is info")
	flag.Parse()

	if *cmd == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	runit.LogLevel = *loglevel

	runner, err := runit.New(*cmd, *watchPath, *alive, *wait)
	if err != nil {
		log.Fatal(err)
	}

	exitStatus, err := runner.Do()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
