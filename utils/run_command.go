package utils

import (
	"bufio"
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

	// Set up pipes for real-time output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command start failed: %v", err)
	}

	// Stream output in real-time
	if Verbose {
		go streamOutput(stdout, os.Stdout)
		go streamOutput(stderr, os.Stderr)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}

func streamOutput(src io.Reader, dest io.Writer) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		fmt.Fprintln(dest, scanner.Text())
	}
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

	// Set up pipes for real-time output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command start failed: %v", err)
	}

	// Stream output in real-time
	if Verbose {
		go streamOutput(stdout, os.Stdout)
		go streamOutput(stderr, os.Stderr)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}
