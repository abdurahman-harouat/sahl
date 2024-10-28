package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

var Verbose bool

func RunCommand(command string) error {
    if Verbose {
        fmt.Printf("\033[1;34m=>\033[0m Running: %s\n", command)
    }

    // Use shell to execute the command to properly expand variables
    cmd := exec.Command("sh", "-c", command)
    
    // Inherit parent environment
    cmd.Env = os.Environ()
    
    // Connect to terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
}

func RunCommandWithSudo(command string) error {
    if Verbose {
        fmt.Printf("\033[1;34m=>\033[0m Running with sudo: %s\n", command)
    }

    var cmd *exec.Cmd
    currentUser, err := user.Current()
    if err != nil {
        return fmt.Errorf("failed to get current user: %v", err)
    }

    if currentUser.Uid != "0" {
        // Use -E flag to preserve environment variables
        cmd = exec.Command("sudo", "-E", "sh", "-c", command)
    } else {
        cmd = exec.Command("sh", "-c", command)
    }

    // Inherit parent environment
    cmd.Env = os.Environ()
    
    // Connect to terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
}