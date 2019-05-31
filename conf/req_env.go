package conf

import (
	"os"
	"strings"
)

// NewReqEnv ...
func NewReqEnv(f *Formation) *ReqEnv {
	names := make([]string, 0)

	for _, p := range f.PublicConfs {
		for rn := range p.Required() {
			names = append(names, strings.ToUpper(rn))
		}
	}

	hasValuesForNames := make(map[string]bool)
	valuesForNames := make(map[string]string)

	for _, envName := range names {
		envValue, ok := os.LookupEnv(envName)
		hasValuesForNames[envName] = ok
		if ok {
			valuesForNames[envName] = envValue
		} else {
			valuesForNames[envName] = envValue
		}
	}

	return &ReqEnv{
		HasValuesForNames: hasValuesForNames,
		ValuesForNames:    valuesForNames,
	}
}

// ReqEnv ...
type ReqEnv struct {
	HasValuesForNames map[string]bool
	ValuesForNames    map[string]string
}
