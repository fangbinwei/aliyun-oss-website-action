package utils

import (
	"path"
	"strings"
)

// Match path with one of patterns
func Match(patterns []string, ossPath string) bool {
	for _, p := range patterns {
		if match(p, ossPath) {
			return true
		}
	}
	return false
}

func match(pattern string, ossPath string) bool {
	pattern = strings.TrimPrefix(pattern, "./")
	if hasMeta(pattern) {
		match, err := path.Match(pattern, ossPath)
		if err != nil {
			return false
		}
		return match
	}
	if !strings.HasPrefix(ossPath, pattern) {
		return false
	}

	// dir
	if strings.HasSuffix(pattern, "/") {
		return true
	}
	// file
	return ossPath == pattern
}

func hasMeta(p string) bool {
	return strings.IndexAny(p, "*?[") >= 0
}
