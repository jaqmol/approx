package axenvs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jaqmol/approx/axmsg"
)

// Envs ...
type Envs struct {
	processorName string
	Required      map[string]string
	Optional      map[string]string
	Ins           []string
	Outs          []string
}

// NewEnvs ...
func NewEnvs(processorName string, requiredEnvs []string, optionalEnvs []string) *Envs {
	required, missing := readAllEnvs(requiredEnvs)
	logFatalIfNeeded(processorName, missing)
	optional, _ := readAllEnvs(optionalEnvs)
	ins, missingIns := readIndexedEnvs(processorName, "IN")
	ins = amendIfNeeded(ins, missingIns, "stdin")
	outs, missingOuts := readIndexedEnvs(processorName, "OUT")
	outs = amendIfNeeded(outs, missingOuts, "stdout")
	missing = append(missingIns, missingOuts...)
	logFatalIfNeeded(processorName, missing)
	return &Envs{
		processorName: processorName,
		Required:      required,
		Optional:      optional,
		Ins:           ins,
		Outs:          outs,
	}
}

// InsOuts ...
func (e *Envs) InsOuts() ([]bufio.Reader, []bufio.Writer) {
	ins, inErrs := Inputs(e)
	outs, outErrs := Outputs(e)
	errs := append(inErrs, outErrs...)
	if len(errs) > 0 {
		errMsg := &axmsg.Errors{Source: e.processorName}
		for _, e := range errs {
			errMsg.Log(nil, e.Error())
		}
		errMsg.LogFatal(nil, "%v errors attempting to open input/output streams", len(errs))
	}
	return ins, outs
}

func amendIfNeeded(values []string, missing []string, amendment string) []string {
	if len(values) == 0 && len(missing) == 0 {
		return append(values, amendment)
	}
	return values
}

func logFatalIfNeeded(processorName string, missing []string) {
	if len(missing) > 0 {
		errMsg := &axmsg.Errors{Source: processorName}
		envNames := strings.Join(missing, ", ")
		errMsg.LogFatal(nil, "Required envs %v not found", envNames)
	}
}

func readAllEnvs(envNames []string) (results map[string]string, missing []string) {
	results = make(map[string]string)
	missing = make([]string, 0)
	for _, name := range envNames {
		value, ok := os.LookupEnv(name)
		if ok {
			results[name] = value
		} else {
			missing = append(missing, name)
		}
	}
	return
}

func readIndexedEnvs(processorName string, prefix string) (results []string, missing []string) {
	results = make([]string, 0)
	missing = make([]string, 0)

	countName := prefix + "_COUNT"
	countStr, ok := os.LookupEnv(countName)
	if !ok {
		return
	}

	count64, err := strconv.ParseInt(countStr, 10, 32)
	if err != nil {
		errMsg := &axmsg.Errors{Source: processorName}
		errMsg.LogFatal(nil, "Error parsing env %v=%v as number: %v", countName, countStr, err.Error())
	}
	count := int(count64)

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%v_%v", prefix, i)
		value, ok := os.LookupEnv(name)
		if ok {
			results = append(results, value)
		} else {
			missing = append(missing, name)
		}
	}
	return
}
