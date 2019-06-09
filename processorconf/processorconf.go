package processorconf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/jaqmol/approx/errormsg"
)

// ProcessorConf ...
type ProcessorConf struct {
	Envs    map[string]string
	Inputs  []*bufio.Reader
	Outputs []*bufio.Writer
}

// NewProcessorConf ...
func NewProcessorConf(
	processorName string,
	requiredEnvs []string,
) *ProcessorConf {
	errMsg := &errormsg.ErrorMsg{Processor: processorName}
	return &ProcessorConf{
		Envs: readAllEnvs(errMsg, requiredEnvs),
		Inputs: openInputs(
			errMsg,
			readAllPrefixedEnvs(
				errMsg,
				"IN_",
				parseIntEnv(errMsg, "IN_COUNT", 0),
				"stdin",
			),
		),
		Outputs: openOutputs(
			errMsg,
			readAllPrefixedEnvs(
				errMsg,
				"OUT_",
				parseIntEnv(errMsg, "OUT_COUNT", 0),
				"stdout",
			),
		),
	}
}

// OptionalEnv ...
func (p *ProcessorConf) OptionalEnv(name string) (string, bool) {
	return os.LookupEnv(name)
}

func readAllEnvs(
	errMsg *errormsg.ErrorMsg,
	envNames []string,
) map[string]string {
	values := make(map[string]string)
	for _, envName := range envNames {
		values[envName] = readEnv(errMsg, envName)
	}
	return values
}

func readEnv(
	errMsg *errormsg.ErrorMsg,
	envName string,
) string {
	value, ok := os.LookupEnv(envName)
	if !ok {
		errMsg.LogFatal(nil, "Required env %v not found", envName)
	}
	return value
}

func openInputs(
	errMsg *errormsg.ErrorMsg,
	inValues []string,
) []*bufio.Reader {
	inputs := make([]*bufio.Reader, 0)
	for _, name := range inValues {
		if name == "stdin" {
			inputs = append(inputs, bufio.NewReader(os.Stdin))
		} else {
			f, err := os.OpenFile(name, os.O_RDONLY, 0600)
			if err != nil {
				errMsg.LogFatal(nil, "Error opening named pipe %v for reading: %v", name, err.Error())
			}
			inputs = append(inputs, bufio.NewReader(f))
		}
	}
	return inputs
}

func openOutputs(
	errMsg *errormsg.ErrorMsg,
	outValues []string,
) []*bufio.Writer {
	outputs := make([]*bufio.Writer, 0)
	for _, name := range outValues {
		if name == "stdout" {
			outputs = append(outputs, bufio.NewWriter(os.Stdout))
		} else {
			f, err := os.OpenFile(name, os.O_RDWR, 0600)
			if err != nil {
				errMsg.LogFatal(nil, "Error opening named pipe %v for writing: %v", name, err.Error())
			}
			outputs = append(outputs, bufio.NewWriter(f))
		}
	}
	return outputs
}

func readAllPrefixedEnvs(
	errMsg *errormsg.ErrorMsg,
	prefix string,
	count int,
	fallback string,
) []string {
	acc := make([]string, 0)
	if count == 0 {
		acc = append(acc, fallback)
	} else {
		for i := 0; i < count; i++ {
			name := readIndexedEnv(errMsg, prefix, i)
			acc = append(acc, name)
		}
	}
	return acc
}

func parseIntEnv(
	errMsg *errormsg.ErrorMsg,
	name string,
	fallback int,
) int {
	valueStr, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	value64, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		errMsg.LogFatal(nil, "Error parsing int env %v: %v", name, err.Error())
	}
	return int(value64)
}

func readIndexedEnv(
	errMsg *errormsg.ErrorMsg,
	prefix string,
	index int,
) string {
	name := fmt.Sprintf("%v%v", prefix, index)
	value, ok := os.LookupEnv(name)
	if !ok {
		errMsg.LogFatal(nil, "Required env %v not found", name)
	}
	return value
}
