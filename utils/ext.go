package utils

import (
	"path"
	"strings"
)

// IsHTML is used to determine if a file is HTML
func IsHTML(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".html")
}

// IsPDF is used to determine if a file is PDF
func IsPDF(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

// IsImage is used to determine if a file is image
func IsImage(filename string) bool {
	imageExts := []string{
		".png",
		".jpg",
		".jpeg",
		".webp",
		".gif",
		".bmp",
		".tiff",
		".ico",
		".svg",
	}
	return func() bool {
		ext := path.Ext(filename)
		for _, e := range imageExts {
			if e == ext {
				return true
			}
		}
		return false
	}()
}
