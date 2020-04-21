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

func startFork(usrWrName string, usrRdNames []string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", usrWrName)

	os.Remove(usrWrFilepath)
	usrWrFile, err := open(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}

	usrRdFilepaths := make([]string, len(usrRdNames))
	usrRdFiles := make([]*os.File, len(usrRdNames))
	for i, usrRdName := range usrRdNames {
		usrRdFilepath := fmt.Sprintf("%s.rd", usrRdName)
		usrRdFilepaths[i] = usrRdFilepath
		os.Remove(usrRdFilepath)
		usrRdFile, err := open(usrRdFilepath)
		if err != nil {
			log.Fatal("Error making/opening named pipe for reading: ", err)
		}
		usrRdFiles[i] = usrRdFile
	}

	usrRdWriters := make([]io.Writer, len(usrRdFiles))
	for i, usrRdFiles := range usrRdFiles {
		usrRdWriters[i] = usrRdFiles
	}

	log.Printf("Running fork from %s to %s\n", usrWrFilepath, strings.Join(usrRdFilepaths, ", "))
	err = runFork(usrWrFile, usrRdWriters)
	if err != nil {
		log.Fatalln("Error operating fork loop:", err)
	}
}

func runFork(usrWrFile io.Reader, usrRdFiles []io.Writer) error {
	scanner := bufio.NewScanner(usrWrFile)
	scanner.Split(scanMessages)

	for scanner.Scan() {
		msgB64 := scanner.Bytes()

		msg := make([]byte, base64.StdEncoding.DecodedLen(len(msgB64)))
		_, err := base64.StdEncoding.Decode(msg, msgB64)
		if err != nil {
			log.Printf("Error decoding message: %s\n", err)
		}
		log.Printf("FORKING: %s\n", msg)

		msgWithDelim := append(msgB64, delim...)
		for _, usrRdFile := range usrRdFiles {
			usrRdFile.Write(msgWithDelim)
		}
	}

	return scanner.Err()
}
