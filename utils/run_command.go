package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"os/user"
	"strings"
)

var Verbose bool

func RunCommand(command string) error {
    // Split the command into the command name and its arguments
    cmdArgs := strings.Fields(command)
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

    // Get the output pipe
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

    // Read and print the output in real-time if verbose is enabled
    if Verbose {
        go printOutput(stdout)
        go printOutput(stderr)
    }

    // Wait for the command to finish
    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
}

// Helper function to print output line by line
func printOutput(pipe io.Reader) {
    reader := bufio.NewReader(pipe)
    for {
        line, _, err := reader.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Printf("error reading output: %v\n", err)
            break
        }
        fmt.Println(string(line))
    }
}

func RunCommandWithSudo(command string) error {
    cmdArgs := strings.Fields(command)
    var cmd *exec.Cmd
    if currentUser, _ := user.Current(); currentUser.Uid != "0" {
        cmd = exec.Command("sudo", append([]string{"sh", "-c"}, command)...)
    } else {
        cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Errorf("failed to get stdout pipe: %v", err)
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return fmt.Errorf("failed to get stderr pipe: %v", err)
    }

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("command start failed: %v", err)
    }

    if Verbose {
        go printOutput(stdout)
        go printOutput(stderr)
    }

    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
}
