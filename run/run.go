package run

import (
	"log"
	"os"
	"strings"

	"github.com/jaqmol/approx/conf"
	"github.com/jaqmol/approx/flow"
)

// Init ...
func Init() (fl *flow.Flow, err error) {
	fo := conf.ReadFormation()
	re := conf.NewReqEnv(fo)
	exitIfRequirementsAreMissing(re)
	fl = flow.NewFlow(fo)
	// hub, err = NewHub(re, fo)
	return
}

// Run ...
func Run(fl *flow.Flow) <-chan error {
	errChan := make(chan error, 0)
	// for _, publicSource := range hub.PublicProcs {
	// 	sources := []proc.Proc{publicSource}
	// 	for sources != nil && len(sources) > 0 {
	// 		for _, s := range sources {
	// 			s.Start(errChan)
	// 		}
	// 		destinations := utils.CollectOuts(sources...)
	// 		if _, ok := destinations[] utils.ContainsProc(destinations, publicSource) {
	// 			sources = nil
	// 		} else {
	// 			sources = destinations
	// 		}
	// 	}
	// }
	return errChan
}

func exitIfRequirementsAreMissing(re *conf.ReqEnv) {
	shouldExit := false
	allNames := make([]string, 0)
	for name, hasValue := range re.HasValuesForNames {
		allNames = append(allNames, name)
		if !hasValue {
			log.Printf("Please provide environment variable: %v\n", name)
			shouldExit = true
		}
	}
	if shouldExit {
		os.Exit(1)
	} else {
		log.Printf("Using environment variables %v\n", strings.Join(allNames, ", "))
	}
}
