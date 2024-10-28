package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

var Verbose bool

// CommandOutput stores both stdout and stderr
type CommandOutput struct {
    Stdout string
    Stderr string
}

func RunCommand(command string) error {
    cmdArgs := strings.Fields(command)
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

    // Create pipes for both stdout and stderr
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdout pipe: %v", err)
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return fmt.Errorf("failed to get stderr pipe: %v", err)
    }

    // Buffer to store output
    var stdoutBuf, stderrBuf strings.Builder
    var wg sync.WaitGroup

    // Start the command
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("command start failed: %v", err)
    }

    // Read both stdout and stderr concurrently
    wg.Add(2)
    go func() {
        defer wg.Done()
        readOutput(stdout, &stdoutBuf, true) // true for stdout
    }()
    go func() {
        defer wg.Done()
        readOutput(stderr, &stderrBuf, false) // false for stderr
    }()

    // Wait for both goroutines to complete
    wg.Wait()

    // Wait for the command to finish
    err = cmd.Wait()
    if err != nil {
        // If there was an error, show both stdout and stderr
        if stdoutBuf.Len() > 0 {
            fmt.Printf("Command output:\n%s\n", stdoutBuf.String())
        }
        if stderrBuf.Len() > 0 {
            fmt.Printf("Error output:\n%s\n", stderrBuf.String())
        }
        return fmt.Errorf("command failed: %v\nError output: %s", err, stderrBuf.String())
    }

    return nil
}

func readOutput(pipe io.Reader, buffer *strings.Builder, isStdout bool) {
    scanner := bufio.NewScanner(pipe)
    for scanner.Scan() {
        line := scanner.Text()
        if Verbose {
            if isStdout {
                fmt.Printf("➜ %s\n", line)
            } else {
                fmt.Printf("❯ %s\n", line)
            }
        }
        buffer.WriteString(line + "\n")
    }
}

func RunCommandWithSudo(command string) error {
    var cmd *exec.Cmd
    currentUser, err := user.Current()
    if err != nil {
        return fmt.Errorf("failed to get current user: %v", err)
    }

    if currentUser.Uid != "0" {
        cmd = exec.Command("sudo", "sh", "-c", command)
    } else {
        cmdArgs := strings.Fields(command)
        cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
    }

    // Create pipes for both stdout and stderr
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdout pipe: %v", err)
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return fmt.Errorf("failed to get stderr pipe: %v", err)
    }

    // Buffer to store output
    var stdoutBuf, stderrBuf strings.Builder
    var wg sync.WaitGroup

    // Start the command
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("command start failed: %v", err)
    }

    // Read both stdout and stderr concurrently
    wg.Add(2)
    go func() {
        defer wg.Done()
        readOutput(stdout, &stdoutBuf, true)
    }()
    go func() {
        defer wg.Done()
        readOutput(stderr, &stderrBuf, false)
    }()

    // Wait for both goroutines to complete
    wg.Wait()

    // Wait for the command to finish
    err = cmd.Wait()
    if err != nil {
        // If there was an error, show both stdout and stderr
        if stdoutBuf.Len() > 0 {
            fmt.Printf("Command output:\n%s\n", stdoutBuf.String())
        }
        if stderrBuf.Len() > 0 {
            fmt.Printf("Error output:\n%s\n", stderrBuf.String())
        }
        return fmt.Errorf("command failed: %v\nError output: %s", err, stderrBuf.String())
    }

    return nil
}