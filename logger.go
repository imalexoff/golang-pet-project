package main

import (
	"fmt"
	"time"
)

func log(message string) {
	fmt.Printf("[%v] %v\n", time.Now().Format(time.Stamp), message)
}

func logf(format string, a ...interface{}) {
	fmt.Printf("[%v] %v\n", time.Now().Format(time.Stamp), fmt.Sprintf(format, a...))
}
