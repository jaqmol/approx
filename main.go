package main

import (
	"log"

	"github.com/jaqmol/approx/run"
	"github.com/jaqmol/approx/visualize"
)

func main() {
	hub, err := run.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	visualize.Hub(hub)
}
