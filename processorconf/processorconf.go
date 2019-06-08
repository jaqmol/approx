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
func NewProcessorConf(processorName string, requiredEnvs []string) *ProcessorConf {
	return &ProcessorConf{
		Envs: readAllEnvs(processorName, requiredEnvs),
		Inputs: openInputs(
			processorName,
			readAllPrefixedEnvs(
				processorName,
				"IN_",
				parseIntEnv(processorName, "IN_COUNT", 0),
				"stdin",
			),
		),
		Outputs: openOutputs(
			processorName,
			readAllPrefixedEnvs(
				processorName,
				"OUT_",
				parseIntEnv(processorName, "OUT_COUNT", 0),
				"stdout",
			),
		),
	}
}

func readAllEnvs(processorName string, envNames []string) map[string]string {
	values := make(map[string]string)
	for _, envName := range envNames {
		values[envName] = readEnv(processorName, envName)
	}
	return values
}

func readEnv(processorName string, envName string) string {
	value, ok := os.LookupEnv(envName)
	if !ok {
		errormsg.LogFatal(processorName, nil, -1001, "Required env %v not found", envName)
	}
	return value
}

func openInputs(processorName string, inValues []string) []*bufio.Reader {
	inputs := make([]*bufio.Reader, 0)
	for _, name := range inValues {
		if name == "stdin" {
			inputs = append(inputs, bufio.NewReader(os.Stdin))
		} else {
			f, err := os.OpenFile(name, os.O_RDONLY, 0600)
			if err != nil {
				errormsg.LogFatal(processorName, nil, -1002, "Error opening named pipe %v for reading: %v", name, err.Error())
			}
			inputs = append(inputs, bufio.NewReader(f))
		}
	}
	return inputs
}

func openOutputs(processorName string, outValues []string) []*bufio.Writer {
	outputs := make([]*bufio.Writer, 0)
	for _, name := range outValues {
		if name == "stdout" {
			outputs = append(outputs, bufio.NewWriter(os.Stdout))
		} else {
			f, err := os.OpenFile(name, os.O_RDWR, 0600)
			if err != nil {
				errormsg.LogFatal(processorName, nil, -1003, "Error opening named pipe %v for writing: %v", name, err.Error())
			}
			outputs = append(outputs, bufio.NewWriter(f))
		}
	}
	return outputs
}

func readAllPrefixedEnvs(processorName string, prefix string, count int, fallback string) []string {
	acc := make([]string, 0)
	if count == 0 {
		acc = append(acc, fallback)
	} else {
		for i := 0; i < count; i++ {
			name := readIndexedEnv(processorName, prefix, i)
			acc = append(acc, name)
		}
	}
	return acc
}

func parseIntEnv(processorName string, name string, fallback int) int {
	valueStr, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	value64, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		errormsg.LogFatal(processorName, nil, -1004, "Error parsing int env %v: %v", name, err.Error())
	}
	return int(value64)
}

func readIndexedEnv(processorName string, prefix string, index int) string {
	name := fmt.Sprintf("%v%v", prefix, index)
	value, ok := os.LookupEnv(name)
	if !ok {
		errormsg.LogFatal(processorName, nil, -1005, "Required env %v not found", name)
	}
	return value
}
