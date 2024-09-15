package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/enescakir/emoji"
	"github.com/schollz/progressbar/v3"
)

// DownloadFile downloads a file from the given URL and displays a progress bar
func DownloadFile(url string) ([]byte, error) {
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
		Timeout:   180 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, response.Status)
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

	// Create a buffer to store the downloaded data
	var buffer []byte
	reader := io.TeeReader(response.Body, bar)

	// Read the response body and update the progress bar
	buffer, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%v Error reading file: %v", emoji.RedCircle, err)
	}

	return buffer, nil
}