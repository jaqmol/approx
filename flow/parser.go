package flow

import (
	"regexp"
	"strings"

	"github.com/jaqmol/approx/utils"
)

// Parse ...
func Parse(formation map[interface{}]interface{}) (procFlow map[string][]string, tappedPipes map[string]string) {
	procFlow = make(map[string][]string)
	tappedPipes = make(map[string]string)
	flowDef := findFlowDefinition(formation)
	splitter := regexp.MustCompile("-(\\w+)->|->")
	for _, lineDefinition := range flowDef {
		procNames := trimmed(splitter.Split(lineDefinition, -1))
		collectProcFlow(procFlow, procNames)
		pipeNameSubs := splitter.FindAllStringSubmatch(lineDefinition, -1)
		pipeNames := trimmed(findPipeNames(pipeNameSubs))
		collectTappedPipes(tappedPipes, procNames, pipeNames)
	}
	return
}

func findFlowDefinition(formation map[interface{}]interface{}) []string {
	for key, value := range formation {
		KEY := strings.ToUpper(key.(string))
		if KEY == "FLOW" {
			interfaceSlice := value.([]interface{})
			definition := make([]string, len(interfaceSlice))
			for i, interfaceLine := range interfaceSlice {
				definition[i] = interfaceLine.(string)
			}
			return definition
		}
	}
	return nil
}

func collectProcFlow(acc map[string][]string, procNames []string) {
	var fromName string
	for i, toName := range procNames {
		if i > 0 {
			tos, ok := acc[fromName]
			if !ok {
				tos = make([]string, 0)
			}
			tos = append(tos, toName)
			acc[fromName] = tos
		}
		fromName = toName
	}
}

func findPipeNames(pipeNameSubs [][]string) []string {
	acc := make([]string, 0)
	for _, nameSub := range pipeNameSubs {
		pipeName := nameSub[len(nameSub)-1]
		acc = append(acc, pipeName)
	}
	return acc
}

func trimmed(procNames []string) []string {
	for i, pn := range procNames {
		procNames[i] = strings.TrimSpace(pn)
	}
	return procNames
}

func collectTappedPipes(acc map[string]string, procNames []string, pipeNames []string) {
	var fromName string
	for i, toName := range procNames {
		if i > 0 {
			pipeName := pipeNames[i-1]
			if pipeName != "" {
				key := utils.PipeKey(fromName, toName)
				acc[key] = pipeName
			}
		}
		fromName = toName
	}
}
