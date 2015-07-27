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

func main() {
	cmd := flag.String("cmd", "", "command to run *required")
	alive := flag.Bool("alive", false, "try to keep the command alive if it dies *optional")
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
	err = runner.Run(*alive)
	if err != nil {
		log.Fatal(err)
	}

	exitChan := make(chan int, 1)
	if !*alive && *watchPath == "" {
		go func() {
			exitStatus := 0
			if err := runner.Wait(); err != nil {
				exitStatus, err = getExitStatus(err)
				if err != nil {
					log.Fatal(err)
				}
			}
			exitChan <- exitStatus
		}()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)
	for {
		select {
		case sig := <-interrupt:
			log.Printf("captured %v", sig)
			switch sig {
			case syscall.SIGHUP:
				runner.Restart()
				continue
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
				os.Exit(0)
			default:
				continue
			}
		case exitCode := <-exitChan:
			os.Exit(exitCode)
		}
	}
}
