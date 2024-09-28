package utils

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"github.com/schollz/progressbar/v3"
)

func DownloadAndSaveFile(url string) (string, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   200 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, response.Status)
	}

	// Try to get the filename from the Content-Disposition header
	var fileName string
	contentDisposition := response.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			fileName = params["filename"]
		}
	}

	// If no filename in header, derive it from the URL
	if fileName == "" {
		urlParts := strings.Split(url, "/")
		fileName = urlParts[len(urlParts)-1]
	}

	// Create a progress bar
	bar := progressbar.NewOptions64(
		response.ContentLength,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(25),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	// Create the output file
	outFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("%v Error creating file: %v", emoji.RedCircle, err)
	}
	defer outFile.Close()

	// Download the file and write it to disk with progress bar
	_, err = io.Copy(io.MultiWriter(outFile, bar), response.Body)
	if err != nil {
		return "", fmt.Errorf("%v Error saving file: %v", emoji.RedCircle, err)
	}

	return fileName, nil
}