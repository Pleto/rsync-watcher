# Rsync-Watcher

## Overview
`rsync-watcher` is a command-line tool that monitors file system changes in a source directory and syncs those changes to a destination directory using `rsync`

## Requirements
- Go 1.22.1 or higher
- `rsync` installed on the system
- `ssh-agent` (for remote SSH synchronization)

## Installation
To build the rsync-watcher from source:

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/rsync-watcher.git
    cd rsync-watcher
    ```

2. Build the binary:
    ```bash
    go build -o rsync-watcher .
    ```

3. Move the binary to a directory in your $PATH for easy access:
    ```bash
    sudo mv rsync-watcher /usr/local/bin
    ```

## Usage
```bash
rsync-watcher -src <source-directory> -dest <destination-directory> [-ssh-add <ssh-key-1>, <ssh-key-2>, ...] [-- <rsync-flags>]
```
- `-src`: The source directory to watch.
- `-dest`: The destination directory to sync changes to.
- `-ssh-add`: Optionally, specify one or more SSH private keys to add to `ssh-agent` for remote syncing.
- `--`: Pass custom flags to override the default `rsync` options.

## Examples
1. Basic local sync:
    ```bash
    rsync-watcher -src /path/to/source -dest /path/to/destination
    ```

    This command will monitor changes in /path/to/source and sync them to /path/to/destination using the default rsync options.

2. Remote sync with SSH key:
    ```bash
    rsync-watcher -src /path/to/source -dest user@remote:/path/to/destination -ssh-add /path/to/ssh/key
    ```

    In this example, files are synchronized to a remote server via SSH, and the specified SSH private key is added to the `ssh-agent`.

3. Using custom rsync flags
    ```bash
    rsync-watcher -src /path/to/source -dest /path/to/destination -- -avz --exclude '*.log' --dry-run
    ```

    The custom rsync flags (--exclude '*.log' and --dry-run) will be used in place of the default flags.