package message

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
)

type frameState int

const (
	awaitingStart frameState = iota
	readingData
)

// FrameReader ...
type FrameReader struct {
	state        frameState
	Index        int64
	Length       int64
	lengthBuffer []byte
	reader       io.Reader
}

var lengthPrefixReadLength int

func init() {
	maxIntStr := fmt.Sprintf("%v", math.MaxInt64)
	lengthPrefixReadLength = len(maxIntStr) + 1 // for the :
}

// GenerateFrameReaders ...
func GenerateFrameReaders(reader io.Reader) <-chan *FrameReader {
	fr := FrameReader{
		reader: reader,
	}
	return &fr
}

// Channel ...
func (fr *FrameReader) Read(b []byte) (n int, err error) {
	n, err = fr.reader.Read(b)

	if err != nil {
		if err == io.EOF {
			log.Fatalln("End of input stream")
		}
		return n, err
	}

	if fr.state == awaitingStart {
		fr.lengthBuffer = append(fr.lengthBuffer, b...)
		colonIndex := bytes.IndexByte(fr.lengthBuffer, ':')
		if colonIndex == -1 {
			if len(fr.lengthBuffer) > lengthPrefixReadLength {
				log.Fatalf("No envelope length found in %v\n", string(fr.lengthBuffer))
			}
		} else {
			lenStr := string(fr.lengthBuffer[:colonIndex])
			length, err := strconv.ParseInt(lenStr, 10, 64)
			if err != nil {
				log.Fatalf("Length could not be parsed from: %v\n", lenStr)
			}
			fr.Length = length
			fr.lengthBuffer = nil
			fr.state = readingData
		}
	} else if fr.state == readingData {

	}

	fr.Index += int64(n)

	return
}
