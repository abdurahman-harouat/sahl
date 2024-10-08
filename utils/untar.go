package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
)

func UntarFile(filePath, destDir string) error {
    // Functions for colorized output
    green := color.New(color.FgGreen).SprintFunc()

    var cmd *exec.Cmd
    var packageDirPath string

    // Check if the file is a .zip file
    if strings.HasSuffix(filePath, ".zip") {
        // Generate the directory name by removing the extension
        packageDirName := RemoveExtension(filepath.Base(filePath))
        packageDirPath = filepath.Join(destDir, packageDirName)

        // Create the directory for extraction
        if err := os.MkdirAll(packageDirPath, 0755); err != nil {
            return fmt.Errorf("%v Error creating package directory: %v", emoji.RedCircle, err)
        }

        // Use bsdtar to extract the .zip file into the newly created directory
        cmd = exec.Command("bsdtar", "-xf", filePath, "-C", packageDirPath)
    } else {
        // Use tar to extract the tar file into the destination directory
        cmd = exec.Command("tar", "-xf", filePath, "-C", destDir)
        packageDirPath = destDir
    }

    // Run the command to extract the package
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("%v Error extracting package: %v", emoji.RedCircle, err)
    }

    if !strings.HasSuffix(filePath, ".zip") {
        // Check if the extracted directory contains multiple items
        files, err := os.ReadDir(packageDirPath)
        if err != nil {
            return fmt.Errorf("%v Error reading extracted directory: %v", emoji.RedCircle, err)
        }

        if len(files) > 1 {
            // Create a new directory to contain the extracted items
            newDir := filepath.Join(destDir, RemoveExtension(filepath.Base(filePath)))
            if err := os.MkdirAll(newDir, 0755); err != nil {
                return fmt.Errorf("%v Error creating new directory: %v", emoji.RedCircle, err)
            }

            // Move the items to the new directory
            for _, file := range files {
                oldPath := filepath.Join(packageDirPath, file.Name())
                newPath := filepath.Join(newDir, file.Name())
                if err := os.Rename(oldPath, newPath); err != nil {
                    return fmt.Errorf("%v Error moving file %s: %v", emoji.RedCircle, file.Name(), err)
                }
            }
            packageDirPath = newDir // Update packageDirPath to the new directory
        }
    }

    if Verbose {
        fmt.Printf("%v Package extracted successfully to %s\n", green(emoji.CheckMark), packageDirPath)
    }
    return nil
}
