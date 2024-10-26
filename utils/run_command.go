package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var Verbose bool

// RunCommand executes a shell command and handles output based on verbosity
func RunCommand(command string) error {
	// Create a shell command
	cmd := exec.Command("sh", "-c", command)

	// Set up output pipes
	var stdout, stderr strings.Builder
	if Verbose {
		// Use MultiWriter to capture output and also display it
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		// Just capture output without displaying it
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Construct detailed error message
		errMsg := fmt.Sprintf("command failed: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg = fmt.Sprintf("command failed with exit code %d", exitErr.ExitCode())
		}

		// Add captured output to error message if there was any
		if stderr.Len() > 0 {
			errMsg += fmt.Sprintf("\nError output:\n%s", stderr.String())
		}
		if stdout.Len() > 0 {
			errMsg += fmt.Sprintf("\nCommand output:\n%s", stdout.String())
		}

		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

// RunCommandWithSudo executes a shell command with sudo if necessary
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

	// Set up output pipes
	var stdout, stderr strings.Builder
	if Verbose {
		// Use MultiWriter to capture output and also display it
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		// Just capture output without displaying it
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	// Run the command
	err = cmd.Run()
	if err != nil {
		// Construct detailed error message
		errMsg := fmt.Sprintf("command failed: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg = fmt.Sprintf("command failed with exit code %d", exitErr.ExitCode())
		}

		// Add captured output to error message if there was any
		if stderr.Len() > 0 {
			errMsg += fmt.Sprintf("\nError output:\n%s", stderr.String())
		}
		if stdout.Len() > 0 {
			errMsg += fmt.Sprintf("\nCommand output:\n%s", stdout.String())
		}

		return fmt.Errorf("%s", errMsg)
	}

	return nil
}