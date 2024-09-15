package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/enescakir/emoji"
)

func UninstallPackage(packageName string) error {
    installPath := fmt.Sprintf("/var/lib/pie/local/%s", packageName)

    // Check if we're running as root
    currentUser, err := user.Current()
    if err != nil {
        return fmt.Errorf("error getting current user: %v", err)
    }
    isRoot := currentUser.Uid == "0"

    // Read installed files
    files, err := ReadLines(filepath.Join(installPath, "installedFiles.txt"), isRoot)
    if err != nil {
        return err
    }

    var errors []string

    // Remove only the files
    for _, file := range files {
        if _, err := os.Stat(file); os.IsNotExist(err) {
            fmt.Printf("Warning: File %s does not exist, skipping\n", file)
            continue
        }

        var cmd *exec.Cmd
        if isRoot {
            cmd = exec.Command("rm", "-f", file)
        } else {
            cmd = exec.Command("sudo", "rm", "-f", file)
        }

        if output, err := cmd.CombinedOutput(); err != nil {
            errors = append(errors, fmt.Sprintf("Could not remove file %s: %v\nOutput: %s", file, err, output))
        }
    }

    // Check and remove empty directories
    dirsToCheck := make(map[string]bool)
    for _, file := range files {
        dir := filepath.Dir(file)
        dirsToCheck[dir] = true
    }

    for dir := range dirsToCheck {
        if isEmpty, _ := IsDirEmpty(dir); isEmpty {
            var cmd *exec.Cmd
            if isRoot {
                cmd = exec.Command("rmdir", dir)
            } else {
                cmd = exec.Command("sudo", "rmdir", dir)
            }

            if output, err := cmd.CombinedOutput(); err != nil {
                errors = append(errors, fmt.Sprintf("Could not remove empty directory %s: %v\nOutput: %s", dir, err, output))
            }
        }
    }

    // Remove the package info directory
    var cmd *exec.Cmd
    if isRoot {
        cmd = exec.Command("rm", "-rf", installPath)
    } else {
        cmd = exec.Command("sudo", "rm", "-rf", installPath)
    }

    if output, err := cmd.CombinedOutput(); err != nil {
        errors = append(errors, fmt.Sprintf("Error removing package info directory: %v\nOutput: %s", err, output))
    }

    if len(errors) > 0 {
        fmt.Printf("%s %s uninstalled with some errors:\n", emoji.Warning, packageName)
        for _, err := range errors {
            fmt.Println(err)
        }
        return fmt.Errorf("uninstallation completed with errors")
    }

    fmt.Printf("%s %s uninstalled successfully\n", emoji.Wastebasket, packageName)
    return nil
}