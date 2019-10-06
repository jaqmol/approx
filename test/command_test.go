package testpackage

import (
	"bytes"
	"log"
	"testing"

	"github.com/jaqmol/approx/message"
	"github.com/jaqmol/approx/processor"

	"github.com/jaqmol/approx/configuration"
)

// TestCommand ...
func TestCommand(t *testing.T) {
	// TODO: NOT WORKING
	originals := loadTestData()[:10]
	originalForID := makePersonForIDMap(originals)
	originalBytes := marshallPeople(originals)

	originalCombined := bytes.Join(originalBytes, configuration.MsgEndBytes)
	originalCombined = append(originalCombined, configuration.MsgEndBytes...)

	reader := bytes.NewReader(originalCombined)
	conf := configuration.Command{
		Ident: "test-command",
		Path:  "/usr/bin/node",
		Args:  []string{"command.js"},
	}
	command := processor.NewCommand(&conf, reader)

	totalCount := 0
	outputReader := command.Outs()[0]
	log.Println("About to start command ...")
	command.Start()
	log.Println("... command did start.")
	scanner := message.NewScanner(outputReader)

	for scanner.Scan() {
		raw := scanner.Bytes()
		data := bytes.Trim(raw, "\x00")
		parsed := checkTestSet(t, originalForID, data)

		log.Printf("Parsed via command.sh: %v\n", parsed)

		totalCount++
	}
}
