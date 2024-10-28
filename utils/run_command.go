package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var Verbose bool

func RunCommand(command string) error {
    cmdArgs := strings.Fields(command)
    cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
    
    // Connect command's stdout and stderr directly to the terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    // Also connect stdin in case the command needs interactive input
    cmd.Stdin = os.Stdin

    if Verbose {
        fmt.Printf("\033[1;34m=>\033[0m Running: %s\n", command)
    }

    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
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

    // Connect command's stdout and stderr directly to the terminal
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    // Also connect stdin in case the command needs interactive input
    cmd.Stdin = os.Stdin

    if Verbose {
        fmt.Printf("\033[1;34m=>\033[0m Running with sudo: %s\n", command)
    }

    err = cmd.Run()
    if err != nil {
        return fmt.Errorf("command failed: %v", err)
    }

    return nil
}