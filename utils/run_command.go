package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

var Verbose bool

// captureOutput reads from a pipe and writes to both a buffer and os.Stdout/os.Stderr
func captureOutput(pipe io.ReadCloser, isStderr bool) (string, error) {
	var output string
	scanner := bufio.NewScanner(pipe)
	
	for scanner.Scan() {
		line := scanner.Text()
		output += line + "\n"
		if Verbose {
			if isStderr {
				fmt.Fprintln(os.Stderr, line)
			} else {
				fmt.Fprintln(os.Stdout, line)
			}
		}
	}
	
	return output, scanner.Err()
}

func RunCommand(command string) error {
	// Create a shell command
	cmd := exec.Command("sh", "-c", command)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Use WaitGroup to ensure we capture all output
	var wg sync.WaitGroup
	var stdoutStr, stderrStr string
	var stdoutErr, stderrErr error

	// Capture stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		stdoutStr, stdoutErr = captureOutput(stdout, false)
	}()

	// Capture stderr
	wg.Add(1)
	go func() {
		defer wg.Done()
		stderrStr, stderrErr = captureOutput(stderr, true)
	}()

	// Wait for output capturing to complete
	wg.Wait()

	// Check for errors in output capturing
	if stdoutErr != nil {
		return fmt.Errorf("error capturing stdout: %v", stdoutErr)
	}
	if stderrErr != nil {
		return fmt.Errorf("error capturing stderr: %v", stderrErr)
	}

	// Wait for command to complete
	err = cmd.Wait()
	if err != nil {
		errMsg := fmt.Sprintf("command failed: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg = fmt.Sprintf("command failed with exit code %d", exitErr.ExitCode())
		}
		
		// Include captured output in error message
		if stderrStr != "" {
			errMsg += fmt.Sprintf("\nError output:\n%s", stderrStr)
		}
		if stdoutStr != "" {
			errMsg += fmt.Sprintf("\nCommand output:\n%s", stdoutStr)
		}
		
		return fmt.Errorf("%s", errMsg)
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

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	// Use WaitGroup to ensure we capture all output
	var wg sync.WaitGroup
	var stdoutStr, stderrStr string
	var stdoutErr, stderrErr error

	// Capture stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		stdoutStr, stdoutErr = captureOutput(stdout, false)
	}()

	// Capture stderr
	wg.Add(1)
	go func() {
		defer wg.Done()
		stderrStr, stderrErr = captureOutput(stderr, true)
	}()

	// Wait for output capturing to complete
	wg.Wait()

	// Check for errors in output capturing
	if stdoutErr != nil {
		return fmt.Errorf("error capturing stdout: %v", stdoutErr)
	}
	if stderrErr != nil {
		return fmt.Errorf("error capturing stderr: %v", stderrErr)
	}

	// Wait for command to complete
	err = cmd.Wait()
	if err != nil {
		errMsg := fmt.Sprintf("command failed: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg = fmt.Sprintf("command failed with exit code %d", exitErr.ExitCode())
		}
		
		// Include captured output in error message
		if stderrStr != "" {
			errMsg += fmt.Sprintf("\nError output:\n%s", stderrStr)
		}
		if stdoutStr != "" {
			errMsg += fmt.Sprintf("\nCommand output:\n%s", stdoutStr)
		}
		
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}