package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jaqmol/approx/configuration"
	"github.com/jaqmol/approx/logger"
	"github.com/jaqmol/approx/message"
)

func main() {
	input := fmt.Sprintf("Hello%vWorld%v", configuration.MessageEnd, configuration.MessageEnd)
	scn := message.NewScanner(strings.NewReader(input))
	for scn.Scan() {
		log.Println(scn.Text())
	}

	l := logger.NewLogger(os.Stdout)
	logMsgs := fmt.Sprintf("LOG%vMESSAGE%v", configuration.MessageEnd, configuration.MessageEnd)
	l.Add(strings.NewReader(logMsgs))
	l.Start()
}
