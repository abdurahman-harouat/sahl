package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
)

func GetOrDownloadPackage(url, cacheDir, expectedMD5 string) (string, error) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

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
				fmt.Printf("%v Package archive already exists and MD5 matches\n", green(emoji.CheckMark))
			} else {
				fmt.Printf("%v Package archive already exists but MD5 does not match\n", yellow(emoji.Warning))
			}
		} else {
			needsDownload = false
			fmt.Printf("%v Package archive already exists, skipping MD5 check\n", green(emoji.CheckMark))
		}
	}

	if needsDownload {
		// Use DownloadAndSaveFile to download the package and save it to disk
		fileName, err := DownloadAndSaveFile(url)  // This will save the file to disk and return the file name
		if err != nil {
			return "", err
		}
	
		// Update filePath to the newly downloaded file
		filePath = fileName
	
		if expectedMD5 != "" {
			calculatedMD5, err := calculateMD5(filePath)
			if err != nil {
				return "", err
			}
			if calculatedMD5 != expectedMD5 {
				return "", fmt.Errorf("%v MD5 verification failed: expected %s, got %s", emoji.RedCircle, expectedMD5, yellow(calculatedMD5))
			}
		}
	
		if Verbose {
			fmt.Printf("%v Package archive downloaded successfully: %s\n", green(emoji.CheckMark), fileName)
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