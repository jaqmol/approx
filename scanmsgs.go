package main

import "bytes"

var delim = []byte{'\n', '-', '-', '-', '\n'}

func scanMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	for i := 0; i+len(delim) <= len(data); {
		j := i + bytes.IndexByte(data[i:], delim[0])
		if j < i {
			break
		}
		if bytes.Equal(data[j+1:j+len(delim)], delim[1:]) {
			// We have a full delim-terminated line.
			return j + len(delim), data[0:j], nil
		}
		i = j + 1
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
