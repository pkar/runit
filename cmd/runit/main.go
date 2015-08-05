package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/pkar/runit"
)

func getExitStatus(waitError error) (int, error) {
	exitError, ok := waitError.(*exec.ExitError)
	if ok {
		waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
		if ok {
			return waitStatus.ExitStatus(), nil
		}
	}
	return 1, waitError
}

func wait(runner *runit.Runner, interrupt chan os.Signal) int {
	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			switch sig {
			case syscall.SIGHUP:
				log.Printf("captured %v restarting...", sig)
				runner.Restart()
				continue
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				log.Printf("captured %v", sig)
				return 0
			default:
				log.Printf("captured %v continue...", sig)
				continue
			}
		}
	}
}

func main() {
	cmd := flag.String("cmd", "", "command to run *required")
	restart := flag.Bool("restart", false, "try to keep the command alive if it dies *optional")
	watchPath := flag.String("watch", "", "path to directory or file to watch *optional")
	flag.Parse()

	if *cmd == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	runner, err := runit.New(*cmd, *watchPath)
	if err != nil {
		log.Fatal(err)
	}
	err = runner.Run(*restart)
	if err != nil {
		log.Fatal(err)
	}

	if !*restart && *watchPath == "" {
		exitStatus := 0
		if err := runner.Wait(); err != nil {
			exitStatus, err = getExitStatus(err)
			if err != nil {
				log.Println(err)
				os.Exit(exitStatus)
			}
		}
		os.Exit(exitStatus)
	}

	interrupt := make(chan os.Signal, 1)
	status := wait(runner, interrupt)
	os.Exit(status)
}
