package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func run(w io.Writer, errw io.Writer) error {
	cfg, err := NewConfig()
	if err != nil {
		return fmt.Errorf("failed to initialise config: %w", err)
	}

	addSSHKeys(w, cfg.SSHAdd)

	src := string(cfg.Source)
	dest := string(cfg.Destination)
	rsyncOverrideFlags := cfg.RSyncOverrideFlags

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		var lastEventTime time.Time
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Multiple events are triggered for a single file write operation.
				// To avoid abusively calling rsync, we'll debounce events, i.e. wait
				// a short period to allow all events related to a single file write to occur,
				// and process them as a group.
				if event.Op&fsnotify.Write == fsnotify.Write {
					now := time.Now()
					if now.Sub(lastEventTime) > 100*time.Millisecond {
						fmt.Fprintln(w, "File system change detected:", event)
						syncDirectories(src, dest, rsyncOverrideFlags)
					}
					lastEventTime = now
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintln(errw, err)
			}
		}
	}()

	err = watcher.Add(src)
	if err != nil {
		return err
	}

	// block the main goroutine, waiting for file system events
	select {}
}

func main() {
	if err := run(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
