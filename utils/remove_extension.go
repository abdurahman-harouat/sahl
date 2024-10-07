package utils

import "strings"

// Function to remove the extension from a package
func RemoveExtension(fileName string) string {
	extensions := []string{".tar.gz", ".tar.xz", ".zip", ".tar.bz2", ".tgz", ".pcf.gz"}
	for _, ext := range extensions {
		if strings.HasSuffix(fileName, ext) {
			return strings.TrimSuffix(fileName, ext)
		}
	}
	return fileName
}