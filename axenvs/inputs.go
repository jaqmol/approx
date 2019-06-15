package axenvs

import (
	"bufio"
	"os"
)

// Inputs ...
func Inputs(envs *Envs) ([]bufio.Reader, []error) {
	inputs := make([]bufio.Reader, 0)
	errors := make([]error, 0)
	for _, name := range envs.Ins {
		if name == "stdin" {
			in := bufio.NewReader(os.Stdin)
			inputs = append(inputs, *in)
		} else {
			f, err := os.OpenFile(name, os.O_RDONLY, 0600)
			if err != nil {
				errors = append(errors, err)
			}
			in := bufio.NewReader(f)
			inputs = append(inputs, *in)
		}
	}
	return inputs, errors
}
