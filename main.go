package main

import (
	"log"

	"github.com/jaqmol/approx/run"
	"github.com/jaqmol/approx/visualize"
)

func main() {
	hub, err := run.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	visualize.Hub(hub)

	// errChan := run.Run(hub)

	// for err := range errChan {
	// 	log.Fatalln(err.Error())
	// }
}
