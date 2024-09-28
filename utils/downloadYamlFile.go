package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/enescakir/emoji"
)

// DownloadFile downloads a file from the given URL and returns its contents as a byte slice (for YAML or text files)
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
		Timeout:   200 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v Error downloading file: %v", emoji.RedCircle, response.Status)
	}

	// Read the response body into memory and return the data
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%v Error reading file contents: %v", emoji.RedCircle, err)
	}

	return data, nil
}