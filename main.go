package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jaqmol/approx/axmsg"
	"github.com/jaqmol/approx/flow"
	"github.com/jaqmol/approx/run"
	"github.com/jaqmol/approx/visualize"
)

func main() {
	errMsg := &axmsg.Errors{Source: "approx"}
	fl := flow.Init(errMsg)
	visualize.Flow(fl)

	rnnr := run.NewRunner(errMsg, fl)
	rnnr.InitProcessors()
	rnnr.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// go func() {
	<-c
	rnnr.Cleanup()
	os.Exit(1)
	// }()

	// for err := range errChan {
	// 	log.Fatalln(err.Error())
	// }
}
