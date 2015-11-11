package main

import (
	"flag"
	"log"
	"os"

	"github.com/pkar/runit"
)

func main() {
	cmd := flag.String("cmd", "", "command to run *required")
	alive := flag.Bool("alive", false, "try to keep the command alive if it dies *optional")
	watchPath := flag.String("watch", "", "path to directory or file to watch and restart cmd *optional")
	loglevel := flag.Int("loglevel", 1, "logging level 1 is info")
	flag.Parse()

	if *cmd == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	runit.LogLevel = *loglevel

	runner, err := runit.New(*cmd, *watchPath, *alive)
	if err != nil {
		log.Fatal(err)
	}

	exitStatus, err := runner.Do()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
