package utils

import (
	"bytes"
	"regexp"
	"sort"
)

const DEFAULT_HEADERS_CONFIG = `[
	{"path": "\\.html$", "headers": {"Cache-Control": "public, max-age=0, must-revalidate"}}
]`

type HeadersConfig struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

func MatchHeadersConfig(path string, headers []HeadersConfig) map[string]string {
	for _, headerConfig := range headers {
		matched, err := regexp.MatchString(headerConfig.Path, path)
		if err != nil {
			continue
		}
		if matched {
			return headerConfig.Headers
		}
	}
	return nil
}

func GetHeadersConfigMD5(headers map[string]string) (string, error) {
	if len(headers) == 0 {
		return "", nil
	}

	// Sort the keys to ensure consistent ordering
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for _, k := range keys {
		buffer.WriteString(k)
		buffer.WriteString(":")
		buffer.WriteString(headers[k])
	}
	return HashMD5(buffer.Bytes())
}
