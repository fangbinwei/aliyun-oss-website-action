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

// LogErrors is used to print []error
func LogErrors(errs []error) {
	if errs != nil {
		fmt.Println("Errors:")
		for i, err := range errs {
			fmt.Printf("%d\n%v\n", i, err)
		}
	}
}
