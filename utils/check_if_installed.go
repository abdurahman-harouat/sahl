package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func CheckIfInstalled(packageName string) bool {
	packagesLog := "/var/log/packages.log"

	file, err := os.Open(packagesLog)
	if err != nil {
		fmt.Println("Error opening logs:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false

	// Compile a regular expression to match the exact package name
	// This assumes package names are listed at the start of each line
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(packageName) + `\s`)

	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading logs:", err)
		return false
	}

	return found
}