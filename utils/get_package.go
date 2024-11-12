package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

const (
	statusSuccess = "[SUCCESS]"
	statusWarning = "[WARNING]"
	statusError   = "[ERROR]"
)

func GetOrDownloadPackage(url, cacheDir, expectedMD5 string) (string, error) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fileName := path.Base(url)
	filePath := filepath.Join(cacheDir, fileName)

	needsDownload := true
	if _, err := os.Stat(filePath); err == nil {
		if expectedMD5 != "" {
			// File exists, check its MD5
			calculatedMD5, err := calculateMD5(filePath)
			if err != nil {
				return "", err
			}
			if calculatedMD5 == expectedMD5 {
				needsDownload = false
				fmt.Printf("%s Package archive already exists and MD5 matches\n", green(statusSuccess))
			} else {
				fmt.Printf("%s Package archive exists but MD5 does not match\n", yellow(statusWarning))
			}
		} else {
			needsDownload = false
			fmt.Printf("%s Package archive exists, skipping MD5 check\n", green(statusSuccess))
		}
	}

	if needsDownload {
		var err error
		filePath, err = DownloadAndSavePackage(url, cacheDir)
		if err != nil {
			return "", fmt.Errorf("%s Failed to download package: %v", red(statusError), err)
		}

		if expectedMD5 != "" {
			calculatedMD5, err := calculateMD5(filePath)
			if err != nil {
				return "", fmt.Errorf("%s Failed to calculate MD5: %v", red(statusError), err)
			}
			if calculatedMD5 != expectedMD5 {
				return "", fmt.Errorf("%s MD5 verification failed: expected %s, got %s", 
					red(statusError), expectedMD5, yellow(calculatedMD5))
			}
		}

		if Verbose {
			fmt.Printf("%s Package downloaded successfully: %s\n", 
				green(statusSuccess), filePath)
		}
	}

	return filePath, nil
}

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error calculating MD5: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}