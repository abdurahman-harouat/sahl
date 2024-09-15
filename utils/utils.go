package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Function to check if a directory is empty
func IsDirEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}


// function to write lines to a file
func WriteLinesToFile(filePath string, lines []string, isRoot bool) error {
	var cmd *exec.Cmd
	content := strings.Join(lines, "\n")

	if isRoot {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cat > %s", filePath))
	} else {
		cmd = exec.Command("sudo", "sh", "-c", fmt.Sprintf("cat > %s", filePath))
	}

	cmd.Stdin = strings.NewReader(content)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error writing to file: %v\nOutput: %s", err, output)
	}

	return nil
}


func ExtractPackage(cachedFilePath string) (string, error) {
    // Get the current working directory
    currentDir, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("failed to get current directory: %w", err)
    }

    // Get the base name of the package file and remove the extension
    baseName := filepath.Base(cachedFilePath)
    extractDirName := RemoveExtension(baseName)

    // Create the full path for the extraction directory
    extractDir := filepath.Join(currentDir, extractDirName)

    // Create the extraction directory
    if err := os.MkdirAll(extractDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create extraction directory: %w", err)
    }

    // Extract the package
    cmd := exec.Command("tar", "-xf", cachedFilePath, "-C", extractDir)
    if err := cmd.Run(); err != nil {
        os.RemoveAll(extractDir)
        return "", fmt.Errorf("failed to extract package: %w", err)
    }

    return extractDir, nil
}