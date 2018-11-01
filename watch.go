package runit

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// Watch watches the runner watch path for changes and
// notifies the runner of change events. The param ignore
// is a list of glob patterns that will not be included in watch.
func (r *Runner) Watch(shutdown <-chan struct{}, ignore []string) (chan bool, error) {
	restartChan := make(chan bool)

	watcher, err := NewRecursiveWatcher(r.WatchPath, ignore)
	if err != nil {
		log.Println("ERRO:", err)
		return nil, err
	}

	go func(restart chan<- bool) {
		for {
			select {
			case <-shutdown:
				watcher.Close()
			case event := <-watcher.Events:
				//log.Println("DEBU: event:", event)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					// create a file or directory
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						log.Println("ERRO:", err)
						continue
					}

					if watcher.ShouldIgnoreFile(event.Name) {
						continue
					}
					if fi.IsDir() {
						//log.Printf("DEBU: detected new directory %s\n", event.Name)
						watcher.AddFolder(event.Name)
						restart <- true
						log.Printf("INFO: added new folder: %s\n", event.Name)
					} else {
						// created a file
						restart <- true
						log.Printf("INFO: added new file: %s\n", event.Name)
					}
				case event.Op&fsnotify.Write == fsnotify.Write:
					if !watcher.ShouldIgnoreFile(event.Name) {
						restart <- true
						log.Printf("INFO: modified file: %s\n", event.Name)
					}
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					watcher.Remove(event.Name)
					restart <- true
					log.Printf("INFO: removed file: %s\n", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Printf("chmod file: %s", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					// renaming a file triggers a create event
					//log.Printf("DEBU: renamed file: %s", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("ERRO:", err)
			}
		}
	}(restartChan)

	return restartChan, nil
}
