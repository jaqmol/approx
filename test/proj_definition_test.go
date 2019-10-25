package test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestProjectDefinition ...
func TestProjectDefinition(t *testing.T) {
	projDir, err := filepath.Abs("simpl-test-proj") // /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	defs, err := project.LoadDefinition(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, defs)
}

func checkProjectDefinitions(t *testing.T, defs map[string]project.Definition) {
	if len(defs) != 4 {
		t.Fatalf("Expected 4 definitions, but got \"%v\"", len(defs))
	}

	forkExp := defs["fork"]
	if forkExp.Ident() != "fork" || forkExp.Type() != project.ForkType {
		t.Fatalf("Expected \"fork\", but got \"%v\"", forkExp.Ident())
	}

	cmd1Exp := defs["extract-first-name"]
	if cmd1Exp.Ident() != "extract-first-name" || cmd1Exp.Type() != project.CommandType {
		t.Fatalf("Expected \"extract-first-name\", but got \"%v\"", cmd1Exp.Ident())
	}
	cmd1 := cmd1Exp.(*project.Command)
	if cmd1.Cmd() != "node ../node-procs/test-extract-prop.js" {
		t.Fatalf("Command 1 cmd mismatch")
	}
	if cmd1.Env()["PROP_NAME"] != "first_name" {
		t.Fatalf("Expected command 1 env PROP_NAME to be \"first_name\", but got \"%v\"", cmd1.Env()["PROP_NAME"])
	}

	cmd2Exp := defs["extract-last-name"]
	if cmd2Exp.Ident() != "extract-last-name" || cmd2Exp.Type() != project.CommandType {
		t.Fatalf("Expected \"extract-last-name\", but got \"%v\"", cmd2Exp.Ident())
	}
	cmd2 := cmd2Exp.(*project.Command)
	if cmd2.Cmd() != "node ../node-procs/test-extract-prop.js" {
		t.Fatalf("Command 2 cmd mismatch")
	}
	if cmd2.Env()["PROP_NAME"] != "last_name" {
		t.Fatalf("Expected command 2 env PROP_NAME to be \"last_name\", but got \"%v\"", cmd2.Env()["PROP_NAME"])
	}

	mergeExp := defs["merge"]
	if mergeExp.Ident() != "merge" || mergeExp.Type() != project.MergeType {
		t.Fatalf("Expected \"merge\", but got \"%v\"", mergeExp.Ident())
	}
}
