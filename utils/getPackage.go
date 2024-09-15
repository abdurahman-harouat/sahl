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
		data, err := DownloadFile(url)
		if err != nil {
			return "", err
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			return "", fmt.Errorf("%v Error writing package archive to file: %v", emoji.RedCircle, err)
		}

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
			fmt.Printf("%v Package archive downloaded successfully\n", green(emoji.CheckMark))
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