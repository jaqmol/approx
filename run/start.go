package run

import (
	"github.com/jaqmol/approx/processor"
)

// Start ...
func Start(procs []processor.Processor) {
	for _, p := range procs {
		p.Start()
	}
}
