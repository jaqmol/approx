package message

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
)

// FrameState ...
type FrameState int

// FrameStates ...
const (
	AwaitingStart FrameState = iota
	ReadingEnvelope
)

// FrameReader ...
type FrameReader struct {
	State          FrameState
	isReading      bool
	expectedLength int64
	ChunkRead      FrameChunkReader
}

// StartReading ...
func (fr *FrameReader) StartReading(reader io.Reader) {
	if fr.isReading {
		log.Fatalln("StartReading must only be called once")
	}
	go fr.startReading(reader)
}

func (fr *FrameReader) startReading(reader io.Reader) {
	readLen := 19
	readBuff := make([]byte, readLen)
	buffer := make([]byte, 0)
	hasBufferedData := false
	var readCount int64
	readCount = 0
	for {
		n, err := reader.Read(readBuff)
		if err == io.EOF {
			log.Fatalln("End of input stream")
		}
		if fr.State == AwaitingStart {
			buffer = append(buffer, readBuff[:n]...)
			hasBufferedData = true
			colonIndex := bytes.IndexByte(buffer, ':')
			if colonIndex == -1 {
				if len(buffer) > readLen {
					log.Fatalf("No envelope length found in %v\n", string(buffer))
				}
			} else {
				lenStr := string(buffer[:colonIndex])
				length, err := strconv.ParseInt(lenStr, 10, 64)
				if err != nil {
					log.Fatalf("Length could not be parsed from: %v\n", lenStr)
				}
				fr.expectedLength = length
				buffer = buffer[colonIndex+1:]
				fr.State = ReadingEnvelope
			}
		} else if fr.State == ReadingEnvelope {
			// TODO: slicing length must be in range of fr.expectedLength
			if readCount < fr.expectedLength {
				if hasBufferedData {
					fr.ChunkRead(fr.expectedLength, readCount, buffer)
					readCount += int64(len(buffer))
					hasBufferedData = false
				}
				data := readBuff[:n]
				readCount += int64(len(data))
				fr.ChunkRead(fr.expectedLength, readCount, data)
			}
		}
		fmt.Print(string(readBuff[:n]))
	}
}

// FrameChunkReader ...
type FrameChunkReader func(length, readCount int64, chunk []byte)

// // Frame ...
// type Frame struct {
// 	reader io.Reader
// 	writer io.Writer
// }
