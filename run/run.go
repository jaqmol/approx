package run

import (
	"log"
	"os"
	"strings"

	"github.com/jaqmol/approx/conf"
)

// Run ...
func Run() (hub *Hub, err error) {
	fo := conf.ReadFormation()
	re := conf.NewReqEnv(fo)
	exitIfRequirementsAreMissing(re)
	hub, err = NewHub(re, fo)
	return
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
