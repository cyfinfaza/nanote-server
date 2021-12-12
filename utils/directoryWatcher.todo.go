package utils

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func WatchDir(path string, callback func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("fsnotify watcher error:", err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					callback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("fsnotify event error:", err)
				}
			}
		}
	}()
	err = watcher.Add(path)
	if err != nil {
		fmt.Println("fsnotify post-add error:", err)
	}
}
