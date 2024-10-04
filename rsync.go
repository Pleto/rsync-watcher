package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func syncDirectories(src, dest, rsyncOverrideFlags string) {
	var cmd *exec.Cmd

	if rsyncOverrideFlags != "" {
		flags := strings.Fields(rsyncOverrideFlags)
		flags = append(flags, src, dest)
		cmd = exec.Command("rsync", flags...)
	} else {
		cmd = exec.Command("rsync", "-avz", "--delete", src, dest)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Rsync error: %s\n", err)
	}
	fmt.Printf("Rsync output: %s\n", output)
}
