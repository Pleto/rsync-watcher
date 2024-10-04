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

    if err := addSSHKeys(w, cfg.SSHAdd); err != nil {
        return fmt.Errorf("addSSHKeys(): %w", err)
    }

	src := string(cfg.Source)
	dest := string(cfg.Destination)
	rsyncOverrideFlags := cfg.RSyncOverrideFlags

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		debouncer := NewEventDebouncer(100*time.Millisecond, func() {
			fmt.Fprintln(w, "Syncing from src to dest")
			if err := syncDirectories(src, dest, rsyncOverrideFlags); err != nil {
				fmt.Fprintln(errw, err)
			}
		})
		debouncer.Trigger()

		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}

				debouncer.Trigger()
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

	// block the main goroutine, waiting for file system events.
	select {}
}

func main() {
	if err := run(os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
