package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func startMerge(usrWrNames []string, usrRdName string) {
	usrWrFilepaths := make([]string, len(usrWrNames))
	usrWrFiles := make([]*os.File, len(usrWrNames))
	for i, usrWrName := range usrWrNames {
		usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)
		usrWrFilepaths[i] = usrWrFilepath
		os.Remove(usrWrFilepath)
		usrWrFile, err := open(usrWrFilepath)
		if err != nil {
			log.Fatal("Error making/opening named pipe for writing: ", err)
		}
		usrWrFiles[i] = usrWrFile
	}

	usrWrReaders := make([]io.Reader, len(usrWrFiles))
	for i, usrWrFiles := range usrWrFiles {
		usrWrReaders[i] = usrWrFiles
	}

	usrRdFilepath := fmt.Sprintf("%s.rd", usrRdName)

	os.Remove(usrRdFilepath)
	usrRdFile, err := open(usrRdFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for reading: ", err)
	}

	log.Printf("Running merge from %s to %s\n", strings.Join(usrWrFilepaths, ", "), usrRdFilepath)
	runMerge(usrWrReaders, usrRdFile)
}

func runMerge(usrWrFiles []io.Reader, usrRdFile io.Writer) {
	merger := make(chan []byte)
	logging := make(chan string)
	errors := make(chan error)
	quit := make(chan error)
	quitCount := 0

	for _, usrWrFile := range usrWrFiles {
		go scanReader(usrWrFile, merger, logging, errors, quit)
	}

	for {
		select {
		case msgWithDelim := <-merger:
			usrRdFile.Write(msgWithDelim)
		case logMsg := <-logging:
			log.Print(logMsg)
		case err := <-errors:
			log.Println(err)
		case err := <-quit:
			if err != nil {
				log.Fatalln("Error operating merge loop:", err)
			}
			quitCount++
			if quitCount == len(usrWrFiles) {
				return
			}
		}
	}
}

func scanReader(
	usrWrFile io.Reader,
	merger chan<- []byte,
	logging chan<- string,
	errors chan<- error,
	quit chan<- error,
) {
	scanner := bufio.NewScanner(usrWrFile)
	scanner.Split(scanMessages)

	for scanner.Scan() {
		msgB64 := scanner.Bytes()

		msg := make([]byte, base64.StdEncoding.DecodedLen(len(msgB64)))
		_, err := base64.StdEncoding.Decode(msg, msgB64)
		if err != nil {
			errMsg := fmt.Errorf("Error decoding message: %s", err)
			errors <- errMsg
		}
		logMsg := fmt.Sprintf("MERGING: %s\n", msg)
		logging <- logMsg

		msgWithDelim := append(msgB64, delim...)
		merger <- msgWithDelim
	}

	errors <- scanner.Err()
}
