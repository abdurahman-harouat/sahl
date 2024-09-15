package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
)

// Function to read and print the content of /var/log/packages.log
func PrintInstalledPackages() {

    packagesLog := "/var/log/packages.log"

    // Open the file
    file, err := os.Open(packagesLog)
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer file.Close()

    // Count the number of lines and print the content
    lineCount := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lineCount++
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading file: %v", err)
    }

    color.Yellow("\n %v Total number of installed packages: %d\n", emoji.Package, lineCount)
}