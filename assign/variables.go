package assign

import (
	"log"
	"strings"

	"github.com/jaqmol/approx/definition"
)

func variables(assigns map[string]string, defForName map[string]definition.Definition) map[string]string {
	varVals := make(map[string]string, 0)
	for varName, valuePathStr := range assigns {
		valuePath := strings.Split(valuePathStr, ".")
		if len(valuePath) != 3 {
			log.Fatalf("Expected assign value path in the scheme of <processor-name>.ENV.<env-name>, but got: %v\n", valuePathStr)
		}
		defName := valuePath[0]
		envName := valuePath[2]
		def, ok := defForName[defName]
		if !ok {
			log.Fatalf("Assign value path error resolving processor: %v\n", defName)
		}
		varVals[varName] = *def.Env[envName]
	}
	return varVals
}
