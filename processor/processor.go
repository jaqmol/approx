package processor

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jaqmol/approx/configuration"
)

/*
	A	processor is initialized with
	- it's specific type of configuration
	- the output(s) of it's predecessor(s) as it's input(s)
*/

// Processor ...
type Processor interface {
	Start()
	Conf() configuration.Processor
	Outs() []io.Reader
	Out() io.Reader
	Err() io.Reader
	Connect(inputs ...io.Reader) error
}

func errorIfInvalidConnect(procID string, inputs []io.Reader, hasPreviousInputs bool) error {
	if !readerSliceIsUsable(inputs) {
		return fmt.Errorf("Connect call on %v must contain inputs", procID)
	}
	if hasPreviousInputs {
		return fmt.Errorf("Connect call on %v can only be performed once", procID)
	}
	return nil
}

func readerSliceIsUsable(readers []io.Reader) bool {
	if readers == nil || len(readers) == 0 {
		return false
	}
	for _, r := range readers {
		if r == nil {
			return false
		}
	}
	return true
}

func inApproxDevEnv() bool {
	approxEnv := strings.ToLower(os.Getenv("APPROX_ENV"))
	return approxEnv == "development" || approxEnv == "dev"
}

func compactString(full interface{}, max int) string {
	msgStr := fmt.Sprintf("%v", full)
	if len(msgStr) > max {
		return msgStr[:max]
	}
	return msgStr
}
