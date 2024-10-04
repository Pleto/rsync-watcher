package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func syncDirectories(src, dest, rsyncOverrideFlags string) error {
	var cmd *exec.Cmd

	if rsyncOverrideFlags != "" {
		flags := strings.Fields(rsyncOverrideFlags)
		flags = append(flags, src, dest)
		cmd = exec.Command("rsync", flags...)
	} else {
		cmd = exec.Command("rsync", "-avz", "--delete", src, dest)
	}

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync error: %w", err)
	}
	return nil
}
