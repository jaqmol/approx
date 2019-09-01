package check

import (
	"log"
	"os/exec"
	"strings"

	"github.com/jaqmol/approx/definition"
)

// Check ...
func Check(definitions []definition.Definition, flows map[string][]string) {
	hasOnePublicDefinition(definitions)
	defTypeForName := makeDefTypeForNameMap(definitions)
	forkesHaveCorrectFlow(flows, defTypeForName)
	mergesHaveCorrectFlow(flows, defTypeForName)
	processesHaveCorrectFlow(flows, defTypeForName)
	executablesExist(definitions)
}

func hasOnePublicDefinition(definitions []definition.Definition) {
	counter := 0
	for _, d := range definitions {
		if d.IsPublic() {
			counter++
			if counter > 1 {
				break
			}
		}
	}
	if counter != 1 {
		log.Fatal("Found more than one public definition")
	}
}

func forkesHaveCorrectFlow(flows map[string][]string, defTypeForName map[string]definition.Type) {
	insCount, outsCount := insAndOutsForType(definition.TypeFork, flows, defTypeForName)

	for name, ins := range insCount {
		if ins > 1 {
			log.Fatalf("Expected fork \"%v\" to have 1 in, but found %v", name, ins)
		}
		outs := outsCount[name]
		if outs < 2 {
			log.Fatalf("Expected fork \"%v\" to have at least 2 outs, but found %v", name, outs)
		}
	}

	for name, outs := range outsCount {
		if outs < 2 {
			log.Fatalf("Expected fork \"%v\" to have at least 2 outs, but found %v", name, outs)
		}
		ins := insCount[name]
		if ins > 1 {
			log.Fatalf("Expected fork \"%v\" to have 1 in, but found %v", name, ins)
		}
	}
}

func mergesHaveCorrectFlow(flows map[string][]string, defTypeForName map[string]definition.Type) {
	insCount, outsCount := insAndOutsForType(definition.TypeMerge, flows, defTypeForName)

	for name, ins := range insCount {
		if ins < 2 {
			log.Fatalf("Expected fork \"%v\" to have at least 2 ins, but found %v", name, ins)
		}
		outs := outsCount[name]
		if outs > 1 {
			log.Fatalf("Expected fork \"%v\" to have 1 out, but found %v", name, outs)
		}
	}

	for name, outs := range outsCount {
		if outs > 1 {
			log.Fatalf("Expected fork \"%v\" to have 1 out, but found %v", name, outs)
		}
		ins := insCount[name]
		if ins < 2 {
			log.Fatalf("Expected fork \"%v\" to have at least 2 ins, but found %v", name, ins)
		}
	}
}

func processesHaveCorrectFlow(flows map[string][]string, defTypeForName map[string]definition.Type) {
	insCount1, outsCount1 := insAndOutsForType(definition.TypeHTTPServer, flows, defTypeForName)
	insCount2, outsCount2 := insAndOutsForType(definition.TypeProcess, flows, defTypeForName)
	insCount := mergeCountMaps(insCount1, insCount2)
	outsCount := mergeCountMaps(outsCount1, outsCount2)

	for name, ins := range insCount {
		if ins > 1 {
			log.Fatalf("Expected process \"%v\" to have 1 in, but found %v", name, ins)
		}
		outs := outsCount[name]
		if outs > 1 {
			log.Fatalf("Expected process \"%v\" to have 1 out, but found %v", name, outs)
		}
	}

	for name, outs := range outsCount {
		if outs > 1 {
			log.Fatalf("Expected process \"%v\" to have 1 out, but found %v", name, outs)
		}
		ins := insCount[name]
		if ins > 1 {
			log.Fatalf("Expected process \"%v\" to have 1 in, but found %v", name, ins)
		}
	}
}

func mergeCountMaps(mapA map[string]int, mapB map[string]int) map[string]int {
	merged := make(map[string]int)
	for k, v := range mapA {
		merged[k] = v
	}
	for k, v := range mapB {
		merged[k] = v
	}
	return merged
}

func makeDefTypeForNameMap(ds []definition.Definition) map[string]definition.Type {
	acc := make(map[string]definition.Type)
	for _, d := range ds {
		acc[d.Name] = d.Type
	}
	return acc
}

func insAndOutsForType(
	defType definition.Type,
	flows map[string][]string,
	defTypeForName map[string]definition.Type,
) (map[string]int, map[string]int) {
	insCount := make(map[string]int)
	outsCount := make(map[string]int)
	for outName, inNames := range flows {
		outType := defTypeForName[outName]
		if defType == outType {
			outs := outsCount[outName]
			outs += len(inNames)
			outsCount[outName] = outs
		}

		for _, inName := range inNames {
			inType := defTypeForName[inName]
			if defType == inType {
				ins := insCount[inName]
				ins++
				insCount[inName] = ins
			}
		}
	}
	return insCount, outsCount
}

func executablesExist(definitions []definition.Definition) {
	notFound := make([]string, 0)
	for _, d := range definitions {
		if d.Type == definition.TypeProcess {
			executable := strings.Split(d.Command, " ")[0]
			_, err := exec.LookPath(executable)
			if err != nil {
				notFound = append(notFound, executable)
			}
		}
	}
	if len(notFound) > 0 {
		log.Fatalf("Process executables not found: %v\n", strings.Join(notFound, ", "))
	}
}
