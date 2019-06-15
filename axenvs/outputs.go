package axenvs

import (
	"bufio"
	"os"
)

// Outputs ...
func Outputs(envs *Envs) ([]bufio.Writer, []error) {
	outputs := make([]bufio.Writer, 0)
	errors := make([]error, 0)
	for _, name := range envs.Outs {
		if name == "stdout" {
			out := bufio.NewWriter(os.Stdout)
			outputs = append(outputs, *out)
		} else {
			f, err := os.OpenFile(name, os.O_RDWR, 0600)
			if err != nil {
				errors = append(errors, err)
			}
			out := bufio.NewWriter(f)
			outputs = append(outputs, *out)
		}
	}
	return outputs, errors
}
