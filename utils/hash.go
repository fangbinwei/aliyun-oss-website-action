package utils

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
)

func HashMD5(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		// TODO: debug info
		return "", err
	}
	defer f.Close()
	return hashMD5(f)

}

func hashMD5(f io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		// TODO: debug info
		return "", err
	}
	result := h.Sum(nil)
	encoded := base64.StdEncoding.EncodeToString(result)

	return encoded, nil
}
