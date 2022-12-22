package main

import (
	"encoding/json"
	"flag"
	"log"
)

type flags struct {
	fd  fileDetails
	uts string //unix to string
}

func main() {

	fg := flags{}

	flag.StringVar(&fg.fd.fpath, "f", "", "--f /home/input.txt (Required)")
	flag.StringVar(&fg.fd.outFile, "o", "", "--f /home/output.txt")
	flag.StringVar(&fg.uts, "uts", "", "--uts createdAt,updatedAt")
	flag.Parse()

	input := fg.fd.getInputReader()
	output := fg.fd.getOutWriter()
	decoder := json.NewDecoder(input)

	csvp := Parser{
		out:     output,
		inp:     input,
		decoder: decoder,
	}

	csvp.setHeadersAndWriteFirstRow().setUTS(fg.uts)

	csvp.parseArray()

	endToken, err := decoder.Token()
	if err != nil {
		log.Fatalf("error while decoding json : %v", err)
	}

	log.Printf("End Token : %v", endToken)

}
