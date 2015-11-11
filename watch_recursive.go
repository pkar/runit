package runit

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-fsnotify/fsnotify"
)

// RecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L13
type RecursiveWatcher struct {
	*fsnotify.Watcher
}

// NewRecursiveWatcher https://github.com/nathany/looper/blob/master/watch.go#L19
func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher}

	pinfo("watching folders in", path)

	rw.AddFolder(path)
	return rw, nil
}

// AddFolder https://github.com/nathany/looper/blob/master/watch.go#L40
func (watcher *RecursiveWatcher) AddFolder(path string) {
	watcher.Add(path)

	subFolders := watcher.Subfolders(path)
	if len(subFolders) == 0 {
		return
	}
	for _, folder := range subFolders {
		err := watcher.Add(folder)
		if err != nil {
			perrorf("error watching: %s %v\n", folder, err)
			return
		}
		pdebugf("added folder %s", folder)
	}
}

// Subfolders returns a slice of subfolders (recursive), including the folder provided.
func (watcher *RecursiveWatcher) Subfolders(path string) (paths []string) {
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

func shouldIgnoreFile(name string) bool {
	return strings.HasPrefix(name, ".")
}
