package test

import "github.com/jaqmol/approx/configuration"

type testProc struct {
	ident string
}

func (f *testProc) Type() configuration.ProcessorType {
	return configuration.ForkType
}

func (f *testProc) ID() string {
	return f.ident
}

// TODO: REMOVE
// func (f *testProc) Next() []configuration.Processor {
// 	return nil
// }

// TODO: REMOVE
// func (f *testProc) SetNext(next ...configuration.Processor) {
// }

func makeTestProcs(count int) []configuration.Processor {
	acc := make([]configuration.Processor, count)
	for i := range acc {
		acc[i] = &testProc{}
	}
	return acc
}
