package run

import (
	"github.com/jaqmol/approx/channel"
	"github.com/jaqmol/approx/definition"
	"github.com/jaqmol/approx/utils"
)

// MakePipes ...
func MakePipes(definitions []definition.Definition, procFlow map[string][]string, tappedPipes map[string]string) map[string]channel.Pipe {
	acc := make(map[string]channel.Pipe)

	for _, fromDef := range definitions {
		fromName := fromDef.Name
		toNames := procFlow[fromName]

		for _, toName := range toNames {
			key := utils.PipeKey(fromName, toName)
			tapName, isTapped := tappedPipes[key]
			if isTapped {
				acc[key] = channel.NewTappedPipe(tapName)
			} else {
				acc[key] = channel.NewPipe()
			}
		}
	}

	return acc
}

// MakeStderrs ...
func MakeStderrs(definitions []definition.Definition, tappedPipes map[string]string) map[string]channel.Pipe {
	acc := make(map[string]channel.Pipe)

	for _, def := range definitions {
		acc[def.Name] = channel.NewPipe()
	}

	for _, name := range tappedPipes {
		acc[name] = channel.NewPipe()
	}

	return acc
}
