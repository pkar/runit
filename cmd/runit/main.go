package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkar/runit"
	"github.com/pkar/runit/vendor/log"
)

func main() {
	logLevel := flag.Uint("loglevel", 1, "logging level, 0 debug *optional")
	cmd := flag.String("cmd", "", "command to run *required")
	watchPath := flag.String("watch", "", "path to directory or file to watch *optional")
	flag.Parse()

	log.SetLevel(*logLevel)
	if *cmd == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	runner, err := runit.Run(*cmd, *watchPath)
	if err != nil {
		log.Fatal(err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			log.Debugf("captured %v", sig)
			switch sig {
			case syscall.SIGHUP:
				runner.Restart()
				continue
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				os.Exit(0)
			default:
				continue
			}
		}
	}
}
