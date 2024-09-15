package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// function to read lines from a file
func ReadLines(filePath string, isRoot bool) ([]string, error) {
	var cmd *exec.Cmd
	if isRoot {
		cmd = exec.Command("cat", filePath)
	} else {
		cmd = exec.Command("sudo", "cat", filePath)
	}
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	
	return strings.Split(string(output), "\n"), nil
}