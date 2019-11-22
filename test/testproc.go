package test

import "github.com/jaqmol/approx/config"

type testProc struct {
	ident string
}

func (f *testProc) Type() config.ProcessorType {
	return config.ForkType
}

func (f *testProc) ID() string {
	return f.ident
}

func makeTestProcs(count int) []config.Processor {
	acc := make([]config.Processor, count)
	for i := range acc {
		acc[i] = &testProc{}
	}
	return acc
}
