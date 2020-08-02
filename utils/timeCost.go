package utils

import (
	"fmt"
	"time"
)
// TimeCost is time-consuming calculation function
// Usage: defer TimeCost()()
func TimeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("time cost = %v\n", tc)
	}
}
