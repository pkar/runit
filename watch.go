package runit

import (
	"os"

	"github.com/go-fsnotify/fsnotify"
	"github.com/pkar/log"
)

// Watch watches the runner watch path for changes and
// notifies the runner of change events. The param ignore
// is a list of glob patterns that will not be included in watch.
func (r *Runner) Watch(shutdown <-chan struct{}, ignore []string) (chan bool, error) {
	restartChan := make(chan bool)

	watcher, err := NewRecursiveWatcher(r.WatchPath, ignore)
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	go func(restart chan<- bool) {
		for {
			select {
			case <-shutdown:
				watcher.Close()
			case event := <-watcher.Events:
				log.Debug.Println("event:", event)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					// create a file or directory
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						log.Error.Println(err)
						continue
					}

					if watcher.ShouldIgnoreFile(event.Name) {
						continue
					}
					if fi.IsDir() {
						log.Debug.Printf("detected new directory %s", event.Name)
						watcher.AddFolder(event.Name)
						restart <- true
						log.Info.Printf("added new folder: %s", event.Name)
					} else {
						// created a file
						restart <- true
						log.Info.Printf("added new file: %s", event.Name)
					}
				case event.Op&fsnotify.Write == fsnotify.Write:
					if !watcher.ShouldIgnoreFile(event.Name) {
						restart <- true
						log.Debug.Printf("modified file: %s", event.Name)
					}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					watcher.Remove(event.Name)
					restart <- true
					log.Info.Printf("removed file: %s", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Debug.Printf("chmod file: %s", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					// renaming a file triggers a create event
					log.Debug.Printf("renamed file: %s", event.Name)
				}
			case err := <-watcher.Errors:
				log.Error.Println(err)
			}
		}
	}(restartChan)

	return restartChan, nil
}
