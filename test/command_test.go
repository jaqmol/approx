package testpackage

import (
	"bytes"
	"testing"

	"github.com/jaqmol/approx/message"
	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"
)

// TestCommand ...
func TestCommand(t *testing.T) {
	originals := loadTestData()
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

	reader := bytes.NewReader(originalCombined)
	conf := configuration.Command{
		Ident: "test-command",
		Cmd:   "node command.js",
	}
	command := processor.NewCommand(&conf, reader)

	totalCount := 0
	outputReader := command.Outs()[0]
	command.Start()
	scanner := message.NewScanner(outputReader)

	for scanner.Scan() {
		raw := scanner.Bytes()
		data := bytes.Trim(raw, "\x00")
		checkTestSet(t, originalForID, data)

		totalCount++
		if totalCount == len(originals) {
			command.SigInt()
		}
	}

	if len(originals) != totalCount {
		t.Fatal("Command dispatch count doesn't corespond to source count")
	}
}

// TestCommandWithLogging ...
func TestCommandWithLogging(t *testing.T) {

}
