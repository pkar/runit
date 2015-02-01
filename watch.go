package runit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkar/runit/vendor/fsnotify"
	"github.com/pkar/runit/vendor/log"
)

// Watch watches the runner watch path for changes and
// notifies the runner of change events.
func (r *Runner) Watch() (chan bool, error) {
	restartChan := make(chan bool)

	watcher, err := NewRecursiveWatcher(r.watchPath)
	if err != nil {
		log.Error(err)
		return restartChan, err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debug("event: ", event)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					// create a file/directory
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						log.Error(err)
						continue
					}

					if fi.IsDir() {
						log.Infof("Detected new directory %s", event.Name)
						if !shouldIgnoreFile(filepath.Base(event.Name)) {
							watcher.AddFolder(event.Name)
						}
						continue
					}
					// created a file
					log.Infof("Detected new file %s", event.Name)
					watcher.Files <- event.Name
					restartChan <- true
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Infof("modified file: %s", event.Name)
					watcher.Files <- event.Name
					restartChan <- true
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Infof("removed file: %s", event.Name)
					restartChan <- true
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Infof("renamed file: %s", event.Name)
					restartChan <- true
				}

			case err := <-watcher.Errors:
				log.Error(err)
			}
		}
	}()

	return restartChan, nil
}

// RecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L13
type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

// NewRecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L19
func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, fmt.Errorf("No folders to watch.")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher}

	rw.Files = make(chan string, 30)
	rw.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

// AddFolder https://github.com/nathany/looper/blob/master/watch.go#L40
func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Errorf("Error watching: %s %v", folder, err)
	}
	watcher.Folders <- folder
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			if shouldIgnoreFile(name) && name != "." && name != ".." {
				return filepath.SkipDir
			}
			paths = append(paths, newPath)
		}
		return nil
	})
	return paths
}

// shouldIgnoreFile determines if a file should be ignored.
// File names that begin with "." or "_" are ignored by the go tool.
func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}
