package processor

// // Stdin ...
// type Stdin struct {
// 	err   *procPipe
// 	inout io.Reader
// }

// // NewStdin ...
// func NewStdin() *Stdin {
// 	return &Stdin{
// 		err:   newProcPipe(),
// 		inout: nil,
// 	}
// }

// // Start ...
// func (s *Stdin) Start() {}

// // Conf ...
// func (s *Stdin) Conf() configuration.Processor {
// 	return configuration.Stdin
// }

// // Outs ...
// func (s *Stdin) Outs() []io.Reader {
// 	return []io.Reader{s.inout}
// }

// // Out ...
// func (s *Stdin) Out() io.Reader {
// 	return s.inout
// }

// // Err ...
// func (s *Stdin) Err() io.Reader {
// 	return s.err.reader()
// }

// // Connect ...
// func (s *Stdin) Connect(inputs ...io.Reader) error {
// 	err := errorIfInvalidConnect(configuration.Stdin.ID(), inputs, s.inout != nil)
// 	if err != nil {
// 		return err
// 	}
// 	s.inout = inputs[0]
// 	return nil
// }