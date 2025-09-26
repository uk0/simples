package services

import (
	"log"
	"os"
	"path/filepath"
	"retroblog/config"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Watch init error: %v", err)
		return
	}
	defer w.Close()

	addTree := func(root string) {
		filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				if err := w.Add(p); err != nil {
					log.Printf("Failed to watch directory %s: %v", p, err)
				}
			}
			return nil
		})
	}

	addTree(config.SrcDir)
	if _, err := os.Stat("tpl"); err == nil {
		addTree("tpl")
	}

	var debounceTimer *time.Timer
	var debounceMu sync.Mutex

	triggerRebuild := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()

		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
			log.Printf("Triggering rebuild...")
			if err := RebuildAll(); err != nil {
				log.Printf("Rebuild error: %v", err)
			} else {
				log.Printf("Rebuild successful")
			}
		})
	}

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}

			if ev.Has(fsnotify.Create) {
				if fi, err := os.Stat(ev.Name); err == nil && fi.IsDir() {
					if err := w.Add(ev.Name); err != nil {
						log.Printf("Failed to watch new directory %s: %v", ev.Name, err)
					}
				}
			}

			if ev.Has(fsnotify.Write) || ev.Has(fsnotify.Create) ||
				ev.Has(fsnotify.Remove) || ev.Has(fsnotify.Rename) {
				log.Printf("Change detected: %s", ev)
				triggerRebuild()
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Printf("Watch error: %v", err)
		}
	}
}
