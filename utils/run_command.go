package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

var Verbose bool

func RunCommand(command string) error {
	// Create a shell command
	cmd := exec.Command("sh", "-c", command)

	// Set the output based on verbosity
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}

func RunCommandWithSudo(command string) error {
	// Check if the user is a sudo user
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	var cmd *exec.Cmd
	if currentUser.Uid != "0" {
		// Run the command with sudo
		cmd = exec.Command("sudo", "sh", "-c", command)
	} else {
		// Run the command without sudo
		cmd = exec.Command("sh", "-c", command)
	}

	// Set the output based on verbosity
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}