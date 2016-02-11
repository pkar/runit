package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkar/runit"
)

func main() {
	cmd := flag.String("cmd", "", `command to run *required`)
	alive := flag.Bool("alive", false, `try to keep the command alive if it dies, you would use this for long running services like a server *optional`)
	watchPath := flag.String("watch", "", `path to directory or file to watch and restart cmd, the command will be run on startup unless wait is specified *optional`)
	ignore := flag.String("ignore", "", `a comma seperated list of regex patterns to ignore, by default hidden files and folders are ignored *optional`)
	wait := flag.Bool("wait", false, `used with watch, this will wait for file changes first and then run the cmd given, rather than run at the beginning *optional`)
	version := flag.Bool("version", false, `print the version and exit`)
	flag.Parse()

	if *version {
		fmt.Println(runit.AppVersion)
		os.Exit(0)
	}

	if *cmd == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	runner, err := runit.New(*cmd, *watchPath, strings.Split(*ignore, ","), *alive, *wait)
	if err != nil {
		log.Fatal(err)
	}

	exitStatus, err := runner.Do()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
