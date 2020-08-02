package utils

import (
	"fmt"
	"os"
)

// HandleError is error handling method, print error and exit
func HandleError(err error) {
	fmt.Println("occurred error:", err)
	os.Exit(1)
}
