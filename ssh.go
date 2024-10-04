package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func addSSHKeys(w io.Writer, sshKeys []string) error {
	if !isSSHAgentRunning() {
		fmt.Fprintln(w, "ssh-agent is not running, starting it now...")
		if err := startSSHAgent(); err != nil {
			return fmt.Errorf("startSSHAgent(): %w", err)
		}
	}

	for _, sshKey := range sshKeys {
		if err := addSSHKey(sshKey); err != nil {
			return fmt.Errorf("addSSHKey: %w", err)
		}
	}

	return nil
}

func isSSHAgentRunning() bool {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	return sshAuthSock != ""
}

func startSSHAgent() error {
	cmd := exec.Command("ssh-agent")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start ssh-agent: %v", err)
	}
	fmt.Printf("ssh-agent started:\n%s\n", output)
	return nil
}

func addSSHKey(privateKeyPath string) error {
	cmd := exec.Command("ssh-add", privateKeyPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add SSH key: %v\n%s", err, output)
	}
	fmt.Printf("SSH key added:\n%s\n", output)
	return nil
}
