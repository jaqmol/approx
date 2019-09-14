package utils

import "fmt"

// PipeKey ...
func PipeKey(fromName string, toName string) string {
	return fmt.Sprintf("%v->%v", fromName, toName)
}
