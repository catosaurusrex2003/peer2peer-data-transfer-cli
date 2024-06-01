package utils

import "fmt"

const red = "\033[31m"
const reset = "\033[0m"

func LogError(format string, args ...interface{}) {
	fmt.Printf(red+format+reset+"\n", args...)
}
