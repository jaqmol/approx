package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jaqmol/approx/errormsg"
	"github.com/jaqmol/approx/flow"
	"github.com/jaqmol/approx/run"
	"github.com/jaqmol/approx/visualize"
)

func main() {
	errMsg := &errormsg.ErrorMsg{Processor: "approx"}
	fl := flow.Init(errMsg)
	visualize.Flow(fl)

	state := run.Flow(errMsg, fl)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	<-c
	state.Cleanup()
	os.Exit(1)
	// }()

	// for err := range errChan {
	// 	log.Fatalln(err.Error())
	// }
}
