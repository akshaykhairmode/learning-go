package main

import (
	"flag"
	"io"
	"log"
	"os"
)

type fileDetails struct {
	fpath   string //the file to read for the json input
	outFile string //the output file path
}

func (f fileDetails) getInputReader() io.ReadCloser {
	if f.fpath == "" {
		flag.PrintDefaults()
		log.Fatal("Input file path cannot be empty")
	}

	fh, err := os.Open(f.fpath)
	if err != nil {
		log.Fatalf("error while opening input file : %v", err)
	}

	return fh
}

func (f fileDetails) getOutWriter() io.WriteCloser {
	if f.outFile == "" {
		return os.Stdout
	}

	fh, err := os.OpenFile(f.outFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("error while creating output file : %v", err)
	}

	return fh
}
