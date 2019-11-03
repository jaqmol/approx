package test

import (
	"log"
	"os"
	"testing"

	"github.com/jaqmol/approx/processor"
)

// TestProcessorFormation ...
func TestProcessorFormation(t *testing.T) {
	// t.SkipNow()
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	testArgs := []string{origArgs[0], "complx-test-proj"}
	os.Args = testArgs

	form, err := processor.NewFormation()
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Continue writing test

	log.Println(form)
}
