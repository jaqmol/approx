package lib

import (
	"io"

	"github.com/jaqmol/procprox/config"
)

// Copy ...
type Copy struct {
	Config config.Copy

	Stdins []io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
