package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
)

func startTap(name string) {
	usrWrFilepath := fmt.Sprintf("%s.wr", name)
	usrRdFilepath := fmt.Sprintf("%s.rd", name)

	os.Remove(usrWrFilepath)
	usrWrFile, err := open(usrWrFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for writing: ", err)
	}
	os.Remove(usrRdFilepath)
	usrRdFile, err := open(usrRdFilepath)
	if err != nil {
		log.Fatal("Error making/opening named pipe for reading: ", err)
	}

	log.Printf("Running tap from %s to %s\n", usrWrFilepath, usrRdFilepath)
	err = runTap(usrWrFile, usrRdFile)
	if err != nil {
		log.Fatalln("Error operating tap loop:", err)
	}
}

func runTap(usrWrFile io.Reader, usrRdFile io.Writer) error {
	scanner := bufio.NewScanner(usrWrFile)
	scanner.Split(scanMessages)

	for scanner.Scan() {
		msgB64 := scanner.Bytes()

		msg := make([]byte, base64.StdEncoding.DecodedLen(len(msgB64)))
		_, err := base64.StdEncoding.Decode(msg, msgB64)
		if err != nil {
			log.Printf("Error decoding message: %s\n", err)
		}
		log.Printf("TAPPING: %s\n", msg)

		msgWithDelim := append(msgB64, delim...)
		usrRdFile.Write(msgWithDelim)
	}

	return scanner.Err()
}
