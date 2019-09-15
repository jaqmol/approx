package channel

import "bytes"

// LineScanner ...
type LineScanner struct {
	reader           Reader
	buffer           []byte
	newline          byte
	lastNewlineIndex int
}

// NewLineScanner ...
func NewLineScanner(reader Reader) *LineScanner {
	ls := LineScanner{
		reader:           reader,
		buffer:           make([]byte, 0),
		newline:          byte('\n'),
		lastNewlineIndex: -1,
	}
	return &ls
}

// Scan ...
func (l *LineScanner) Scan() bool {
	for l.lastNewlineIndex == -1 {
		l.advance()
	}
	return true
}

// // Err ...
// func (l *LineScanner) Err() error {

// }

// Lines ...
func (l *LineScanner) Lines() [][]byte {
	lines := bytes.Split(l.buffer, []byte{l.newline})
	l.buffer = nil
	l.lastNewlineIndex = -1
	return lines
}

func (l *LineScanner) advance() {
	b := <-l.reader.Read()
	l.buffer = append(l.buffer, b...)
	l.lastNewlineIndex = bytes.LastIndexByte(l.buffer, l.newline)
}
